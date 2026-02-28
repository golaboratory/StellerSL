# TaskInputBody


## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**$schema** | **string** | A URL to the JSON Schema for this object. | [optional] [readonly] [default to undefined]
**assigned_to** | **string** |  | [optional] [default to undefined]
**description** | **string** |  | [default to undefined]
**due_date** | **string** |  | [optional] [default to undefined]
**priority** | **number** |  | [default to 0]
**project_id** | **string** |  | [optional] [default to undefined]
**status** | **string** |  | [default to StatusEnum_Todo]
**title** | **string** |  | [default to undefined]

## Example

```typescript
import { TaskInputBody } from './api';

const instance: TaskInputBody = {
    $schema,
    assigned_to,
    description,
    due_date,
    priority,
    project_id,
    status,
    title,
};
```

[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)
