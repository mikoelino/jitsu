import {CodeInTabs, CodeTab} from "../../../components/Code";
import {Hint} from '../../../components/documentationComponents'

# Google Authorization

This page about Google Authorization configuration.
EventNative works with a number of services provided by Google. They all have the same authorization mechanics.

* **Service Account** is a special account with id that look like [NAME@PROJECT.iam.gserviceaccount.com](https://console.cloud.google.com/iam-admin/serviceaccounts/details/107095565645971338726?project=exalted-cogency-279115). 
  The account can have a key \(or several keys\) which is represented by JSON. Please note, to use Service Account as an authorization mechanism, the resource \(google doc, analytics account, add account etc\) should be shared with that account

### Service account configuration

The easiest way to create Service Account is through [**Google Cloud Console**](https://console.cloud.google.com/)**:** go to Navigation ("burger" at top right corner) → IAM & Admin → Service account. Create service account and download key JSON.

There a few other ways (including console utils), please see [documentation](https://cloud.google.com/iam/docs/creating-managing-service-account-keys)

Service Account key can be referred in a few ways in EventNative configuration:

<CodeInTabs>
    <CodeTab title="As a JSON object" lang="yaml">
        {`
        auth:
  service_account_key: {
    "type": "service_account",
    "project_id": "<PROJECT_ID>",
    "private_key_id": "<PK_ID>",
    "private_key": "<PRIVATE_KEY>",
    "client_email": "<EMAIL>",
    "client_id": "CID",
    "auth_uri": "https://accounts.google.com/o/oauth2/auth",
    "token_uri": "https://oauth2.googleapis.com/token",
    "auth_provider_x509_cert_url": "https://www.googleapis.com/oauth2/v1/certs",
    "client_x509_cert_url": "<CERT_URL>"
  }
        `}
    </CodeTab>
    <CodeTab title="As a JSON string" lang="yaml">
        {`
        auth:
          service_account_key: '{"type":"service_account","project_id":"<PROJECT_ID>","private_key_id":"<PK_ID>","private_key":"<PRIVATE_KEY>","client_email":"<EMAIL>","client_id":"CID","auth_uri":"https://accounts.google.com/o/oauth2/auth","token_uri":"https://oauth2.googleapis.com/token","auth_provider_x509_cert_url":"https://www.googleapis.com/oauth2/v1/certs","client_x509_cert_url":"<CERT_URL>"}'
        `}
    </CodeTab>
    <CodeTab title="As a path to file" lang="yaml">
        {`
        auth:
          service_account_key: '/path/to/file.json'
        `}
    </CodeTab>
</CodeInTabs>