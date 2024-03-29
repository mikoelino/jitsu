---
sort: "0001"
title: "JS SDK (new)"

---

import {Hint} from "../../../../components/documentationComponents";
import {LargeLink} from "../../../../components/documentationComponents";

# Jitsu JS SDK

Jitsu JavaScript SDK allows you to send data to Jitsu from web-application. There're two ways
how to embed Jitsu into your application:

 * Using [npm (or yarn) package manager](/docs/sending-data/js-sdk/package). This is a preferred way for applications written from scratch with
a frameworks such as React, Vue, Angular etc.
 * Using [html snippet](/docs/sending-data/js-sdk/snippet) (also known as 'Pixel' or 'Tracking' code). This is a preferred way for websites and
applications built with no-code tools such as Webflow, Wix, Wordpress etc
 * Also it's possible to use Jitsu react hooks to simplify collecting data from React applications. Those hooks
assume that you use react-router
 * Examples and bindings for other popular frameworks are coming

<Hint>
    We have a deprecated version of JS SDK with non-compatible API. If you're looking to migrate to a new API,
    check out <b><a href="/docs/sending-data/js-sdk/migrating">Migration Guide</a></b>
</Hint>



<LargeLink title="NPM/Yarn package" href="/docs/sending-data/js-sdk/package" />
<LargeLink title="HTML snippet guide" href="/docs/sending-data/js-sdk/snippet" />
<LargeLink title="React guide" href="/docs/sending-data/js-sdk/react" />




