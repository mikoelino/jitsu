---
sort: "0002"

---

import {Hint} from "../../../../components/documentationComponents";
import {LargeLink} from "../../../../components/documentationComponents";
import {CodeInTabs, CodeTab} from "../../../../components/Code";

# NPM package

Jitsu JS SDK available as an npm package. Npm/yarn is a preferred way of sending data to Jitsu
for applications build with modern frameworks such as React, Angular, Vue etc.

Of you're using no-code tools such as Webflow, Wix, Wordpress etc, [html snippet](/docs/sending-data/js-sdk/snippet) will work better

<Hint>
    Check out <a href="/docs/sending-data/js-sdk/react">Jitsu React hooks</a>. Special guides for other frameworks are coming soon!
</Hint>

## Installing Jitsu

Use the following command to add it to your project:

<CodeInTabs>
    <CodeTab title="npm" lang="bash">
        npm install --save @jitsu/sdk-js
    </CodeTab>
    <CodeTab title="yarn" lang="javascript">
        yarn add @jitsu/sdk-js
    </CodeTab>
</CodeInTabs>

To initialize **Jitsu**, please use:

```javascript
const { jitsuClient } = require('@jitsu/sdk-js');
const jitsu = jitsuClient({
    key: "[API_KEY]",
    ...params
});
```

<Hint>
    <a href="/docs/sending-data/js-sdk/parameters-reference">Please see the full list of parameters</a>, a <b>key</b> parameter value is required.
</Hint>


## Sending data

Jitsu exposes only two methods `id()` - for identifyling users and `track()` for sending events.

### ID method

`id()` sets the properties of the user (such as `name`, `e-mail`, `internal id` — any parameters are accepted)

```javascript
jitsu.id({
    "name": "Man with No Name",
    "email": "thegoods@western.com",
    "internal_id": "6"
})
```

By default, a `user_identification` event will be sent. However, it can be changed by setting the second parameter to `true`.

```javascript
jitsu.id({...}, true);
```

### Track method

`eventN.track()` is used to record any event that happens on a webpage

<CodeInTabs>
    <CodeTab title="Syntax" lang="javascript">
        {"eventN.track('{event_name}', {...event data})"}
    </CodeTab>
    <CodeTab title="Example" lang="javascript">
        {`
        eventN.track('product_page_view', {
            product_id: '1e48fb70-ef12-4ea9-ab10-fd0b910c49ce',
            product_price: 399,
            price_currency: 'USD'
        });
        `}
    </CodeTab>
</CodeInTabs>

<Hint>
    If <code inline="true">id()</code> has been called prior to <code inline="true">track()</code>, all user details will be included.
    Also, do not include parameters such as page URL and user agent. Jitsu collects this automatically!
    Read more about <a href="/docs/configuration/schema-and-mappings">our event scheme</a>.
</Hint>

## Intercepting Segment events

As Jitsu, can serve as [Segment replacement](/docs/other-features/segment-compatibility), you can optionally
intercept events that has been sent to segment before

Preferred way of doing that would be supplying jitsu with Analytics object explicitely

```javascript
    const jitsu = jitsuClient({
        key: "[API_KEY]",
        ...params
    });

    //Create analytics via npm module
    const analytics = new Analytics();
    //initialize interception explicitly
    jitsu.interceptAnalytics(analytics);
```

However, if analytics.js is injected as code snippet, not as a package following code will do the job:

```javascript
    const jitsu = jitsuClient({
        key: "[API_KEY]",
        segment_hook: true // instruct jitsu to automatically intercept events
    });
```

<Hint>
    Please make sure that this code executed <b>before</b> initialization of Segment's analytics.js
</Hint>




