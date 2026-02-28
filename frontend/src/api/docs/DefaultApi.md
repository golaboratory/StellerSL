# DefaultApi

All URIs are relative to *http://localhost*

|Method | HTTP request | Description|
|------------- | ------------- | -------------|
|[**addTeamMember**](#addteammember) | **POST** /teams/{id}/members | Add Team Member|
|[**bulkCreateTasks**](#bulkcreatetasks) | **POST** /tasks/bulk | Bulk Create Tasks|
|[**bulkDeleteTasks**](#bulkdeletetasks) | **DELETE** /tasks/bulk | Bulk Delete Tasks|
|[**bulkUpdateTasksStatus**](#bulkupdatetasksstatus) | **PATCH** /tasks/bulk/status | Bulk Update Tasks Status|
|[**createProject**](#createproject) | **POST** /projects | Create Project|
|[**createTask**](#createtask) | **POST** /tasks | Create Task|
|[**createTeam**](#createteam) | **POST** /teams | Create Team|
|[**deleteTask**](#deletetask) | **DELETE** /tasks/{id} | Delete Task|
|[**getDashboardStats**](#getdashboardstats) | **GET** /dashboard | Dashboard Stats|
|[**getUserGrowth**](#getusergrowth) | **GET** /growth | User Growth|
|[**listBadges**](#listbadges) | **GET** /badges | User Badges|
|[**listProjects**](#listprojects) | **GET** /projects | List Projects|
|[**listTasks**](#listtasks) | **GET** /tasks | List Tasks|
|[**listTeams**](#listteams) | **GET** /teams | List Teams|
|[**login**](#login) | **POST** /auth/login | User Login|
|[**register**](#register) | **POST** /auth/register | User Registration|
|[**updateTaskStatus**](#updatetaskstatus) | **PATCH** /tasks/{id}/status | Update Task Status|

# **addTeamMember**
> addTeamMember(teamMemberInputBody)


### Example

```typescript
import {
    DefaultApi,
    Configuration,
    TeamMemberInputBody
} from './api';

const configuration = new Configuration();
const apiInstance = new DefaultApi(configuration);

let id: string; // (default to undefined)
let teamMemberInputBody: TeamMemberInputBody; //

const { status, data } = await apiInstance.addTeamMember(
    id,
    teamMemberInputBody
);
```

### Parameters

|Name | Type | Description  | Notes|
|------------- | ------------- | ------------- | -------------|
| **teamMemberInputBody** | **TeamMemberInputBody**|  | |
| **id** | [**string**] |  | defaults to undefined|


### Return type

void (empty response body)

### Authorization

No authorization required

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/problem+json


### HTTP response details
| Status code | Description | Response headers |
|-------------|-------------|------------------|
|**204** | No Content |  -  |
|**0** | Error |  -  |

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **bulkCreateTasks**
> bulkCreateTasks(bulkTaskCreateInputBody)


### Example

```typescript
import {
    DefaultApi,
    Configuration,
    BulkTaskCreateInputBody
} from './api';

const configuration = new Configuration();
const apiInstance = new DefaultApi(configuration);

let bulkTaskCreateInputBody: BulkTaskCreateInputBody; //

const { status, data } = await apiInstance.bulkCreateTasks(
    bulkTaskCreateInputBody
);
```

### Parameters

|Name | Type | Description  | Notes|
|------------- | ------------- | ------------- | -------------|
| **bulkTaskCreateInputBody** | **BulkTaskCreateInputBody**|  | |


### Return type

void (empty response body)

### Authorization

No authorization required

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/problem+json


### HTTP response details
| Status code | Description | Response headers |
|-------------|-------------|------------------|
|**204** | No Content |  -  |
|**0** | Error |  -  |

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **bulkDeleteTasks**
> bulkDeleteTasks(bulkTaskDeleteInputBody)


### Example

```typescript
import {
    DefaultApi,
    Configuration,
    BulkTaskDeleteInputBody
} from './api';

const configuration = new Configuration();
const apiInstance = new DefaultApi(configuration);

let bulkTaskDeleteInputBody: BulkTaskDeleteInputBody; //

const { status, data } = await apiInstance.bulkDeleteTasks(
    bulkTaskDeleteInputBody
);
```

### Parameters

|Name | Type | Description  | Notes|
|------------- | ------------- | ------------- | -------------|
| **bulkTaskDeleteInputBody** | **BulkTaskDeleteInputBody**|  | |


### Return type

void (empty response body)

### Authorization

No authorization required

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/problem+json


### HTTP response details
| Status code | Description | Response headers |
|-------------|-------------|------------------|
|**204** | No Content |  -  |
|**0** | Error |  -  |

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **bulkUpdateTasksStatus**
> bulkUpdateTasksStatus(bulkTaskUpdateInputBody)


### Example

```typescript
import {
    DefaultApi,
    Configuration,
    BulkTaskUpdateInputBody
} from './api';

const configuration = new Configuration();
const apiInstance = new DefaultApi(configuration);

let bulkTaskUpdateInputBody: BulkTaskUpdateInputBody; //

const { status, data } = await apiInstance.bulkUpdateTasksStatus(
    bulkTaskUpdateInputBody
);
```

### Parameters

|Name | Type | Description  | Notes|
|------------- | ------------- | ------------- | -------------|
| **bulkTaskUpdateInputBody** | **BulkTaskUpdateInputBody**|  | |


### Return type

void (empty response body)

### Authorization

No authorization required

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/problem+json


### HTTP response details
| Status code | Description | Response headers |
|-------------|-------------|------------------|
|**204** | No Content |  -  |
|**0** | Error |  -  |

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **createProject**
> ProjectOutputBody createProject(projectInputBody)


### Example

```typescript
import {
    DefaultApi,
    Configuration,
    ProjectInputBody
} from './api';

const configuration = new Configuration();
const apiInstance = new DefaultApi(configuration);

let projectInputBody: ProjectInputBody; //

const { status, data } = await apiInstance.createProject(
    projectInputBody
);
```

### Parameters

|Name | Type | Description  | Notes|
|------------- | ------------- | ------------- | -------------|
| **projectInputBody** | **ProjectInputBody**|  | |


### Return type

**ProjectOutputBody**

### Authorization

No authorization required

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json, application/problem+json


### HTTP response details
| Status code | Description | Response headers |
|-------------|-------------|------------------|
|**200** | OK |  -  |
|**0** | Error |  -  |

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **createTask**
> TaskOutputBody createTask(taskInputBody)


### Example

```typescript
import {
    DefaultApi,
    Configuration,
    TaskInputBody
} from './api';

const configuration = new Configuration();
const apiInstance = new DefaultApi(configuration);

let taskInputBody: TaskInputBody; //

const { status, data } = await apiInstance.createTask(
    taskInputBody
);
```

### Parameters

|Name | Type | Description  | Notes|
|------------- | ------------- | ------------- | -------------|
| **taskInputBody** | **TaskInputBody**|  | |


### Return type

**TaskOutputBody**

### Authorization

No authorization required

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json, application/problem+json


### HTTP response details
| Status code | Description | Response headers |
|-------------|-------------|------------------|
|**200** | OK |  -  |
|**0** | Error |  -  |

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **createTeam**
> createTeam(teamInputBody)


### Example

```typescript
import {
    DefaultApi,
    Configuration,
    TeamInputBody
} from './api';

const configuration = new Configuration();
const apiInstance = new DefaultApi(configuration);

let teamInputBody: TeamInputBody; //

const { status, data } = await apiInstance.createTeam(
    teamInputBody
);
```

### Parameters

|Name | Type | Description  | Notes|
|------------- | ------------- | ------------- | -------------|
| **teamInputBody** | **TeamInputBody**|  | |


### Return type

void (empty response body)

### Authorization

No authorization required

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/problem+json


### HTTP response details
| Status code | Description | Response headers |
|-------------|-------------|------------------|
|**204** | No Content |  * ID -  <br>  * Name -  <br>  |
|**0** | Error |  -  |

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **deleteTask**
> deleteTask()


### Example

```typescript
import {
    DefaultApi,
    Configuration
} from './api';

const configuration = new Configuration();
const apiInstance = new DefaultApi(configuration);

let id: string; // (default to undefined)

const { status, data } = await apiInstance.deleteTask(
    id
);
```

### Parameters

|Name | Type | Description  | Notes|
|------------- | ------------- | ------------- | -------------|
| **id** | [**string**] |  | defaults to undefined|


### Return type

void (empty response body)

### Authorization

No authorization required

### HTTP request headers

 - **Content-Type**: Not defined
 - **Accept**: application/problem+json


### HTTP response details
| Status code | Description | Response headers |
|-------------|-------------|------------------|
|**204** | No Content |  -  |
|**0** | Error |  -  |

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **getDashboardStats**
> DashboardOutputBody getDashboardStats()


### Example

```typescript
import {
    DefaultApi,
    Configuration
} from './api';

const configuration = new Configuration();
const apiInstance = new DefaultApi(configuration);

const { status, data } = await apiInstance.getDashboardStats();
```

### Parameters
This endpoint does not have any parameters.


### Return type

**DashboardOutputBody**

### Authorization

No authorization required

### HTTP request headers

 - **Content-Type**: Not defined
 - **Accept**: application/json, application/problem+json


### HTTP response details
| Status code | Description | Response headers |
|-------------|-------------|------------------|
|**200** | OK |  -  |
|**0** | Error |  -  |

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **getUserGrowth**
> GrowthOutputBody getUserGrowth()


### Example

```typescript
import {
    DefaultApi,
    Configuration
} from './api';

const configuration = new Configuration();
const apiInstance = new DefaultApi(configuration);

const { status, data } = await apiInstance.getUserGrowth();
```

### Parameters
This endpoint does not have any parameters.


### Return type

**GrowthOutputBody**

### Authorization

No authorization required

### HTTP request headers

 - **Content-Type**: Not defined
 - **Accept**: application/json, application/problem+json


### HTTP response details
| Status code | Description | Response headers |
|-------------|-------------|------------------|
|**200** | OK |  -  |
|**0** | Error |  -  |

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **listBadges**
> BadgeListOutputBody listBadges()


### Example

```typescript
import {
    DefaultApi,
    Configuration
} from './api';

const configuration = new Configuration();
const apiInstance = new DefaultApi(configuration);

const { status, data } = await apiInstance.listBadges();
```

### Parameters
This endpoint does not have any parameters.


### Return type

**BadgeListOutputBody**

### Authorization

No authorization required

### HTTP request headers

 - **Content-Type**: Not defined
 - **Accept**: application/json, application/problem+json


### HTTP response details
| Status code | Description | Response headers |
|-------------|-------------|------------------|
|**200** | OK |  -  |
|**0** | Error |  -  |

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **listProjects**
> ProjectListOutputBody listProjects()


### Example

```typescript
import {
    DefaultApi,
    Configuration
} from './api';

const configuration = new Configuration();
const apiInstance = new DefaultApi(configuration);

const { status, data } = await apiInstance.listProjects();
```

### Parameters
This endpoint does not have any parameters.


### Return type

**ProjectListOutputBody**

### Authorization

No authorization required

### HTTP request headers

 - **Content-Type**: Not defined
 - **Accept**: application/json, application/problem+json


### HTTP response details
| Status code | Description | Response headers |
|-------------|-------------|------------------|
|**200** | OK |  -  |
|**0** | Error |  -  |

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **listTasks**
> TaskListOutputBody listTasks()


### Example

```typescript
import {
    DefaultApi,
    Configuration
} from './api';

const configuration = new Configuration();
const apiInstance = new DefaultApi(configuration);

const { status, data } = await apiInstance.listTasks();
```

### Parameters
This endpoint does not have any parameters.


### Return type

**TaskListOutputBody**

### Authorization

No authorization required

### HTTP request headers

 - **Content-Type**: Not defined
 - **Accept**: application/json, application/problem+json


### HTTP response details
| Status code | Description | Response headers |
|-------------|-------------|------------------|
|**200** | OK |  -  |
|**0** | Error |  -  |

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **listTeams**
> TeamListOutputBody listTeams()


### Example

```typescript
import {
    DefaultApi,
    Configuration
} from './api';

const configuration = new Configuration();
const apiInstance = new DefaultApi(configuration);

const { status, data } = await apiInstance.listTeams();
```

### Parameters
This endpoint does not have any parameters.


### Return type

**TeamListOutputBody**

### Authorization

No authorization required

### HTTP request headers

 - **Content-Type**: Not defined
 - **Accept**: application/json, application/problem+json


### HTTP response details
| Status code | Description | Response headers |
|-------------|-------------|------------------|
|**200** | OK |  -  |
|**0** | Error |  -  |

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **login**
> LoginOutputBody login(loginInputBody)


### Example

```typescript
import {
    DefaultApi,
    Configuration,
    LoginInputBody
} from './api';

const configuration = new Configuration();
const apiInstance = new DefaultApi(configuration);

let loginInputBody: LoginInputBody; //

const { status, data } = await apiInstance.login(
    loginInputBody
);
```

### Parameters

|Name | Type | Description  | Notes|
|------------- | ------------- | ------------- | -------------|
| **loginInputBody** | **LoginInputBody**|  | |


### Return type

**LoginOutputBody**

### Authorization

No authorization required

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json, application/problem+json


### HTTP response details
| Status code | Description | Response headers |
|-------------|-------------|------------------|
|**200** | OK |  -  |
|**0** | Error |  -  |

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **register**
> register(registerInputBody)


### Example

```typescript
import {
    DefaultApi,
    Configuration,
    RegisterInputBody
} from './api';

const configuration = new Configuration();
const apiInstance = new DefaultApi(configuration);

let registerInputBody: RegisterInputBody; //

const { status, data } = await apiInstance.register(
    registerInputBody
);
```

### Parameters

|Name | Type | Description  | Notes|
|------------- | ------------- | ------------- | -------------|
| **registerInputBody** | **RegisterInputBody**|  | |


### Return type

void (empty response body)

### Authorization

No authorization required

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/problem+json


### HTTP response details
| Status code | Description | Response headers |
|-------------|-------------|------------------|
|**204** | No Content |  -  |
|**0** | Error |  -  |

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **updateTaskStatus**
> TaskOutputBody updateTaskStatus(taskStatusUpdateInputBody)


### Example

```typescript
import {
    DefaultApi,
    Configuration,
    TaskStatusUpdateInputBody
} from './api';

const configuration = new Configuration();
const apiInstance = new DefaultApi(configuration);

let id: string; // (default to undefined)
let taskStatusUpdateInputBody: TaskStatusUpdateInputBody; //

const { status, data } = await apiInstance.updateTaskStatus(
    id,
    taskStatusUpdateInputBody
);
```

### Parameters

|Name | Type | Description  | Notes|
|------------- | ------------- | ------------- | -------------|
| **taskStatusUpdateInputBody** | **TaskStatusUpdateInputBody**|  | |
| **id** | [**string**] |  | defaults to undefined|


### Return type

**TaskOutputBody**

### Authorization

No authorization required

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json, application/problem+json


### HTTP response details
| Status code | Description | Response headers |
|-------------|-------------|------------------|
|**200** | OK |  -  |
|**0** | Error |  -  |

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

