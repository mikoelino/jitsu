---
sort: "0001"
title_: "JS Reference (DEPRECATED)"

---

import {Hint} from "../../../../components/documentationComponents";

# JavaScript Reference

<Hint className="hidden" style={{display: 'none'}}>
    This section refers to a <a href="https://www.npmjs.com/package/@jitsu/eventnative">deprecated client library</a>. The library
    works and will work in foreseeable future (also, we will fix critical bugs). However, new features will be
    added to a <a href="https://www.npmjs.com/package/@jitsu/sdk-js">new version</a> only.
    <br /><br />
    Please, read <a href="/docs/sending-data/js-sdk">JS SDK</a> section if you're starting new project. If you're migrating existing project,
    check out <a href="/docs/sending-data/js-sdk/migration">Migration Guide</a>
</Hint>



## Understanding event tracking

**EventNative**'s JavaScript snippet can work in two modes:

* `Direct tracking pixel`: You'll need to call EventNative `eventN.track(...)` explicitly to send events to the server.
* `Intercept Mode`: We capture events from 3rd-party systems (GoogleAnalytics or Segment) if you have them installed. After inserting one line of JavaScript, EventNative will automatically intercept events and send them to your desired destination. Your original event pipeline will not be disrupted.

## Quickstart

The quickest way to get started with JavaScript integration is to open the welcome page at your instance of EventNative: `https://[your-instance]/p/welcome.html` and use tracking code builder:

![Welcome page](/img/docs/welcomehtml.png)

Check `Google Analytics Interceptor` or `Segment Interceptor` if you want EventNative to listen to 3rd-party events. Check `Send Data Directly` otherwise. Do not forget to replace `[API_KEY]` with your [authorization](/docs/configuration/authorization#clientserver-secrets-authorization).

<Hint>
    Make sure that the code inserted after GA and Segment code if it's working as an interceptor. Add <code inline={true}>eventN.track('pageview')</code> if you're sending events directly.
</Hint>

### Further reading

* Check [JS Configuration](/docs/sending-data/javascript-reference/initialization-parameters) to see a full list of parameters that can be used in `eventnConfig` variable
* Installation with [npm or yarn](/docs/sending-data/javascript-reference/npm-or-yarn)

