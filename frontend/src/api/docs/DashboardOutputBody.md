# DashboardOutputBody


## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**$schema** | **string** | A URL to the JSON Schema for this object. | [optional] [readonly] [default to undefined]
**completed_tasks** | **number** |  | [default to undefined]
**daily_activity** | [**Array&lt;Item&gt;**](Item.md) |  | [default to undefined]
**pending_tasks** | **number** |  | [default to undefined]
**total_tasks** | **number** |  | [default to undefined]

## Example

```typescript
import { DashboardOutputBody } from './api';

const instance: DashboardOutputBody = {
    $schema,
    completed_tasks,
    daily_activity,
    pending_tasks,
    total_tasks,
};
```

[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)
