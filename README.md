# IRSA Example Repository

本仓库展示了如何使用 InitContainer 和主容器共享同一目录的方式，为主容器初始化 AWS 的 Access Key (AK) 和 Secret Key (SK)。

## 使用说明
### 前置条件
必须配置IRSA，如果想要不依赖IRSA，需要自己初始化火山KMS的SDK
https://www.volcengine.com/docs/6460/1324613
### 环境变量
在使用本仓库时，需要设置以下环境变量：

| 环境变量名       | 必须 | 默认值          | 说明                                   |
|------------------|------|-----------------|----------------------------------------|
| `VOLCENGINE_REGION` | 是  | `cn-beijing`    | 火山引擎的区域                         |
| `KMS_SECRET_NAME`   | 是  | `default-secret`| KMS 中存储的 Secret 名称               |
| `KMS_AK_KEY`        | 是  | `accessKey`     | Secret JSON 中表示 AK 的键名           |
| `KMS_SK_KEY`        | 是  | `secretKey`     | Secret JSON 中表示 SK 的键名           |
| `CLOUD_TYPE`        | 否  | `aws`       | 云服务类型，例如 `aws` 表示 AWS 云服务 |

### 样例 Deployment
以下是一个示例 Kubernetes Deployment，展示了如何使用 InitContainer 初始化 AK 和 SK：

```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: irsa-example
spec:
  replicas: 1
  selector:
    matchLabels:
      app: irsa-example
  template:
    metadata:
      labels:
        app: irsa-example
    spec:
      containers:
      - name: main-container
        image: your-main-container-image
        volumeMounts:
        - name: shared-volume
          mountPath: /shared
        env:
        - name: AWS_SHARED_CREDENTIALS_FILE
          value: /shared/credentials
      initContainers:
      - name: init-container
        image: your-init-container-image
        volumeMounts:
        - name: shared-volume
          mountPath: /shared
        env:
        - name: VOLCENGINE_REGION
          value: "cn-beijing"
        - name: KMS_SECRET_NAME
          value: "your-kms-secret-name"
        - name: KMS_AK_KEY
          value: "accessKey"
        - name: KMS_SK_KEY
          value: "secretKey"
        - name: CLOUD_TYPE
          value: "aws"
      volumes:
      - name: shared-volume
        emptyDir: {}
```

### 工作原理
1. InitContainer 从 KMS 中获取存储的 AK 和 SK。
2. InitContainer 将 AK 和 SK 写入共享目录中的 AWS 凭证文件。
3. 主容器通过环境变量 `AWS_SHARED_CREDENTIALS_FILE` 访问共享目录中的凭证文件。

### 注意事项
- 请确保 InitContainer 和主容器共享的目录具有正确的权限设置。
- 确保 KMS 中的 Secret 格式为 JSON，且包含指定的 AK 和 SK 键名。

如有任何问题，请参考代码或联系维护者。