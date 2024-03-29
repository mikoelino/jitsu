/* eslint-disable */
import { ApplicationConfiguration, FeatureSettings, setDebugInfo } from './ApplicationServices';
import { User } from './model';
// @ts-ignore
import { eventN, Tracker } from '@jitsu/eventnative';
import LogRocket from 'logrocket';
import murmurhash from 'murmurhash';
import posthog from 'posthog-js';
import { isNullOrUndef } from '../commons/utils';

const AnalyticsJS = require('./analyticsjs-wrapper.js').default;

type ConsoleMessageListener = (level: string, ...args) => void;

class ConsoleLogInterceptor {
  private initialized: boolean = false;
  private listeners: ConsoleMessageListener[] = [];
  private originalError: (message?: any, ...optionalParams: any[]) => void;

  public addListener(listener: ConsoleMessageListener) {
    this.listeners.push(listener);
  }

  public init() {
    if (this.initialized) {
      return;
    }
    let interceptor = this;

    (function () {
      interceptor.originalError = console.error;
      console.error = function () {
        interceptor.listeners.forEach((listener: ConsoleMessageListener) => {
          try {
            listener('error', arguments);
          } catch (e) {
            console.warn('Error applying error listener');
          }
        });
        interceptor.originalError.apply(this, Array.prototype.slice.call(arguments));
      };
    })();
  }

  public error(message?: any, ...optionalParams: any[]) {
    this.originalError.apply(this, Array.prototype.slice.call(arguments));
  }
}

type ApiErrorInfo = {
  method: string;
  url: string;
  requestPayload: any;
  responseStatus: number;
  responseObject?: any;
  errorMessage?: string;
};

class ApiErrorWrapper extends Error {
  private apiDetails: ApiErrorInfo;

  constructor(message: string, apiDetails: ApiErrorInfo) {
    super(message);
    this.apiDetails = apiDetails;
  }
}

function findError(args: any): Error {
  if (typeof args === 'string' || typeof args === 'number' || typeof args === 'boolean') {
    return null;
  }
  args = Array.prototype.slice.call(args);
  for (let i = 0; i < args.length; i++) {
    let arg = args[i];
    if (isError(arg)) {
      return arg;
    } else if (Array.isArray(Array.prototype.slice.call(arg))) {
      let error = findError(arg);
      if (error) {
        return error;
      }
    }
  }
  return null;
}

function isError(obj: any) {
  if (isNullOrUndef(obj)) {
    return false;
  }
  return (
    obj instanceof Error || obj.constructor.name === 'Error' || (obj.message !== undefined && obj.stack !== undefined)
  );
}

export default class AnalyticsService {
  private globalErrorListenerPresent: boolean = false;
  private appConfig: ApplicationConfiguration;
  private user: User;
  private logRocketInitialized: boolean = false;
  private consoleInterceptor: ConsoleLogInterceptor = new ConsoleLogInterceptor();
  private _anonymizeUsers = false;
  private _appName = 'saas';

  constructor(appConfig: ApplicationConfiguration) {
    this.appConfig = appConfig;
    this.consoleInterceptor.init();
    if (this.appConfig.rawConfig.keys.eventnative) {
      const cfg = {
        key: this.appConfig.rawConfig.keys.eventnative,
        tracking_host: 'https://t.jitsu.com',
        cookie_domain: 'jitsu.com',
        randomize_url: true
      };
      eventN.init(cfg);
    }
    this.setupGlobalErrorHandler();
    this.consoleInterceptor.addListener((level, ...args) => {
      let error = findError(args);
      if (error) {
        this.onGlobalError(error, true);
      }
    });
  }

  public ensureLogRocketInitialized() {
    if (!this.logRocketInitialized && this.appConfig.rawConfig.keys.logrocket) {
      LogRocket.init(this.appConfig.rawConfig.keys.logrocket);
      setDebugInfo('logRocket', LogRocket, false);
      this.logRocketInitialized = true;
    }
  }

  public userHasDomain(email: string, domains: string[]) {
    return domains.find((domain) => email.indexOf('@' + domain) > 0) !== undefined;
  }

  public onUserKnown(user: User) {
    if (!user) {
      return;
    }
    this.user = user;
    if (this.appConfig.rawConfig.keys.posthog) {
      posthog.init(this.appConfig.rawConfig.keys.posthog, { api_host: this.appConfig.rawConfig.keys.posthog_host });
      posthog.people.set({ email: user.email });
      posthog.identify(user.uid);
    }
    this.ensureLogRocketInitialized();
    LogRocket.identify(user.uid, {
      email: user.email
    });
    if (this.appConfig.rawConfig.keys.ajs) {
      AnalyticsJS.init(this.appConfig.rawConfig.keys.ajs);
      AnalyticsJS.get().identify(user.uid, {
        email: user.email
      });
    }
    if (this.appConfig.rawConfig.keys.eventnative) {
      const payload = this.getJitsuIdPayload(user);
      eventN.id(payload);
    }
  }

  public getJitsuIdPayload({ email, uid }) {
    return {
      email: this._anonymizeUsers ? undefined : email,
      internal_id: this._anonymizeUsers ? 'hid_' + murmurhash.v3(email || uid) : uid
    };
  }

  public withJitsu(callback: (jitsu: Tracker) => void) {
    if (this.appConfig.rawConfig.keys.eventnative) {
      callback(eventN);
    }
  }

  public async withJitsuSync(callback: (jitsu: Tracker) => Promise<void>) {
    if (this.appConfig.rawConfig.keys.eventnative) {
      return await callback(eventN);
    }
    return Promise.resolve();
  }

  public configure(features: FeatureSettings) {
    this._anonymizeUsers = features.anonymizeUsers;
    this._appName = features.appName;
  }

  private isDev() {
    return this.appConfig.appEnvironment === 'development';
  }

  public onPageLoad({ pagePath }: { pagePath: string }) {
    if (this.appConfig.rawConfig.keys.eventnative) {
      eventN.track('app_page', {
        path: pagePath,
        app: this._appName
      });
    }

    if (this.user && this.appConfig.rawConfig.keys.ajs) {
      AnalyticsJS.get().page('app_page', pagePath, {
        app: this._appName
      });
    }

    if (this.user && this.appConfig.rawConfig.keys.posthog) {
      posthog.capture('$pageview');
    }
  }

  public onGlobalError(error: Error, doNotLog?: boolean) {
    if (!doNotLog) {
      //call console log through interceptor, to make sure it won't be handled
      this.consoleInterceptor.error('[Jitsu] uncaught error', error);
    }
    if (!this.isDev()) {
      try {
        this.sendException(error);
      } catch (e) {
        console.warn('Failed to send event to error monitoring', e);
      }
    }
  }

  public onGlobalErrorEvent(event: ErrorEvent) {
    this.consoleInterceptor.error(
      `[Jitsu] uncaught error '${event.message || 'unknown'}' at ${event.filename}:${event.lineno}:${event.colno}`,
      event.error
    );
    if (!this.isDev()) {
      try {
        this.sendException(event.error);
      } catch (e) {
        console.warn('Failed to send event to error monitoring', e);
      }
    }
  }

  setupGlobalErrorHandler() {
    if (!this.globalErrorListenerPresent) {
      window.addEventListener('error', (event) => this.onGlobalErrorEvent(event));
      window.addEventListener('unhandledrejection', (event) => {
        this.onGlobalError(new Error('Unhandled rejection: ' + event.reason));
      });
      this.globalErrorListenerPresent = true;
    }
  }

  onFailedAPI(param: ApiErrorInfo) {
    let message = `[Jitsu] ${param.method.toUpperCase()} ${param.url} failed with ${param.responseStatus}`;
    this.consoleInterceptor.error(message);
    if (!this.isDev()) {
      this.sendException(new ApiErrorWrapper(message, param));
    }
  }

  private sendException(error: Error) {
    if (!this.isDev()) {
      console.log('Sending error to monitoring system');
      this.ensureLogRocketInitialized();
      if (this.appConfig.rawConfig.keys.logrocket) {
        LogRocket.captureException(error, {
          tags: {
            environment: window.location.host
          }
        });
      }
    }
  }
}

declare global {
  interface Window {
    analytics: any;
  }
}
