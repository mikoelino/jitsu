import {APIParam, APIMethod} from "../../../components/documentationComponents";

# Sync Tasks

Jitsu supports automatic collection synchronization as well as manual. For using automatic collection synchronization
there must be configured `schedule` property in the `collection` section of configuration (see above).

For using manual collection synchronization(HTTP API) there must be admin token configuration.

```yaml
server:
  admin_token: your_admin_token
  
sources:
  source_id:
    type: ...
    ...
```
<br/>

<APIMethod method="POST" path="/api/v1/tasks" title="Running sync task"/>

Since there can be only one task per source - collection pair in the task queue, EventNative returns ID of an existing task, or a new one.
(HTTP responses will have different HTTP codes - see example below)
Authorization admin token might be provided either as query parameter or HTTP header.

<h4>Parameters</h4>

<APIParam name={"source"} dataType="string" required={true} type="queryString" description="Source ID from 'sources' configuration section"/>
<APIParam name={"collection"} dataType="string" required={true} type="queryString" description="Collection name from 'sources' configuration section"/>
<APIParam name={"X-Admin-Token"} dataType="string" required={true} type="header" description="Admin token"/>
<APIParam name={"token"} dataType="string" required={true} type="queryString" description="Admin token"/>

<h4>Response</h4>

Task has been created:

```json
HTTP 201 Created

{
    "task_id": "$sourceId_$collectionName_$UUID"
}
```

Task already exists:

```json
HTTP 200 OK

{
    "task_id": "$sourceId_$collectionName_$UUID" #id of an existing task
}
```

<h4>Error Response</h4>

Source wasn't found:

```json
{
    "message": "Error getting source",
    "error": "Source [jitsu_firebase] doesn't exist"
}
```

<h4>Authorization Error Response</h4>

```json
{
    "message": "Admin token does not match"
}
```

<h4> CURL example</h4>

```bash
curl --location --request POST 'https://<your_server>/api/v1/tasks?source=<your_source_id>&collection=<your_collection_name>&token=<admin_token>'
```

<br/>

<APIMethod method="GET" path="/api/v1/tasks" title="Get all sync tasks"/>

Authorization admin token might be provided either as query parameter or HTTP header

<h4>Parameters</h4>

<APIParam name={"source"} dataType="string" required={true} type="queryString" description="Source ID from 'sources' configuration section"/>
<APIParam name={"collection"} dataType="string" required={false} type="queryString" description="Collection name from 'sources' configuration section. Default value: all collections"/>
<APIParam name={"start"} dataType="string" required={true} type="queryString" description="Start of time interval in ISO 8601 ('2006-01-02T15:04:05.000000Z') format" />
<APIParam name={"end"} dataType="string" required={true} type="queryString" description="End of time interval in ISO 8601 ('2006-01-02T15:04:05.000000Z') format" />
<APIParam name={"limit"} dataType="int" required={false} type="queryString" description="Limit of returned tasks per collection. Default value: 0 - no limit" />
<APIParam name={"status"} dataType="string" required={false} type="queryString" description="Task status filter. Available values: [scheduled, running, failed, success]. Default value: all statuses" />
<APIParam name={"X-Admin-Token"} dataType="string" required={true} type="header" description="Admin token"/>
<APIParam name={"token"} dataType="string" required={true} type="queryString" description="Admin token"/>

<h4>Response</h4>

Sync tasks list

```json
{
    "tasks": [
        {
            "id": "$sourceId_$collectionName_$UUID",
            "source": "$sourceId",
            "collection": "$collectionName",
            "priority": 299998384585588,
            "created_at": "2021-03-10T22:13:32.433956Z",
            "started_at": "2021-03-10T22:13:32.567439Z",
            "finished_at": "2021-03-10T22:13:34.116187Z",
            "status": "SUCCESS"
        },
        {
            "id": "$sourceId_$collectionName_$UUID",
            "source": "$sourceId",
            "collection": "$collectionName",
            "priority": 299998384585588,
            "created_at": "2021-03-11T00:13:32.433956Z",
            "started_at": "2021-03-11T00:13:32.567439Z",
            "status": "RUNNING"
        }
    ]
}
```

<h4>Error Response</h4>

Source wasn't found:

```json
{
    "message": "Error getting source",
    "error": "Source [jitsu_firebase_auth_uses] doesn't exist"
}
```

<h4>Authorization Error Response</h4>

```json
{
    "message": "Admin token does not match"
}
```

<h4> CURL example</h4>

```bash
curl -X GET 'https://<your_server>/api/v1/tasks?source=<your_source_id>&token=<admin_token>&start=2020-01-01T00:00:00Z&end=2024-12-31T23:59:59Z'
```

<br/>

<APIMethod method="GET" path="/api/v1/tasks/:taskId" title="Get sync task by ID"/>

Authorization admin token might be provided either as query parameter or HTTP header

<h4>Parameters</h4>

<APIParam name={"taskId"} dataType="string" required={true} type="pathParam" description="Task ID"/>
<APIParam name={"X-Admin-Token"} dataType="string" required={true} type="header" description="Admin token"/>
<APIParam name={"token"} dataType="string" required={true} type="queryString" description="Admin token"/>

<h4>Response</h4>

Sync task payload

```json
{
    "id": "$sourceId_$collectionName_$UUID",
    "source": "$sourceId",
    "collection": "$collectionName",
    "priority": 299998384583699,
    "created_at": "2021-03-10T22:45:01.512528Z",
    "status": "SCHEDULED"
}
```

<h4>Error Response</h4>

Source wasn't found:

```json
{
    "message": "Error getting source",
    "error": "Source [jitsu_firebase_auth_uses] doesn't exist"
}
```

<h4>Authorization Error Response</h4>

```json
{
    "message": "Admin token does not match"
}
```

<h4> CURL example</h4>

```bash
curl -X GET 'https://<your_server>/api/v1/tasks/<your_task_id>?token=<admin_token>'
```

<br/>

<APIMethod method="GET" path="/api/v1/tasks/:taskId/logs" title="Get sync task logs"/>

Authorization admin token might be provided either as query parameter or HTTP header

<h4>Parameters</h4>

<APIParam name={"taskId"} dataType="string" required={true} type="pathParam" description="Task ID"/>
<APIParam name={"start"} dataType="string" required={false} type="queryString" description="Start of time interval in ISO 8601 ('2006-01-02T15:04:05.000000Z') format. Default value: Unix start epoch (1970-01-01..)" />
<APIParam name={"end"} dataType="string" required={false} type="queryString" description="End of time interval in ISO 8601 ('2006-01-02T15:04:05.000000Z') format. Default value: time.Now() UTC" />
<APIParam name={"X-Admin-Token"} dataType="string" required={true} type="header" description="Admin token"/>
<APIParam name={"token"} dataType="string" required={true} type="queryString" description="Admin token"/>

<h4>Response</h4>

Sync task log messages

```json
{
    "logs": [
        {
            "time": "2021-03-10T22:45:02.578999Z",
            "message": "[$sourceId_$collectionName_$UUID] Running task...",
            "level": "info"
        },
        {
            "time": "2021-03-10T22:45:02.588929Z",
            "message": "[$sourceId_$collectionName_$UUID] Total intervals: [1]",
            "level": "info"
        },
        {
            "time": "2021-03-10T22:45:03.870479Z",
            "message": "[$sourceId_$collectionName_$UUID] FINISHED SUCCESSFULLY in [1.28] seconds (~ 0.02 minutes)",
            "level": "info"
        }
    ]
}
```

<h4>Error Response</h4>

Source wasn't found:

```json
{
    "message": "Error getting source",
    "error": "Source [jitsu_firebase_auth_uses] doesn't exist"
}
```

<h4>Authorization Error Response</h4>

```json
{
    "message": "Admin token does not match"
}
```

<h4> CURL example</h4>

```bash
curl -X GET 'https://<your_server>/api/v1/tasks/<your_task_id>/logs?token=<admin_token>'
```