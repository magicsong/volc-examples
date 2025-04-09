# VKE Sglang Kubernetes 部署示例

本项目提供了一个脚本和模板，用于根据 `.env` 文件中的配置生成最终的 Kubernetes YAML 文件，该文件能够在VKE上一键拉起一个SGLANG的大模型服务

## 文件结构

- `inference/all-in-one.yaml`  
  包含 Kubernetes 资源的模板文件，使用占位符（如 `${NAMESPACE}`）表示需要替换的变量。

- `scripts/generate_yaml.sh`  
  脚本文件，读取 `.env` 文件中的配置并替换模板中的变量，生成最终的 YAML 文件。

- `.env`  
  配置文件（唯一需要用户提供的文件， `touch .env`）。

## `.env` 文件变量说明

以下是 `.env` 文件中所有变量的说明：
> 已经是最小化配置，所有参数都需要填写

| 变量名            | 描述                                                                                     | 示例值                                                                                     |
|-------------------|------------------------------------------------------------------------------------------|------------------------------------------------------------------------------------------|
| `NAMESPACE`       | Kubernetes 命名空间，用于指定资源所在的命名空间。 请确保namespace存在                                        | `dpsk`                                                                                   |
| `IMAGE`           | Sglang容器镜像地址，用于指定部署的容器镜像。 由于镜像太大，VKE无法提供官方镜像，需要用户自己拉取                                                  | `lmsysorg/sglang:latest`              |
| `GPU`             | GPU 数量，用于指定容器所需的 GPU 资源。                                                  | `1`                                                                                      |
| `BUCKET`          | 对象存储桶名称，用于指定存储模型文件的存储桶。                                           | `demo-bucket-hll`                                                                        |
| `MODEL_PATH_TOS`  | 模型路径，用于指定存储桶中的模型文件路径。                                               | `/models/`                                                                               |
| `URL`             | 对象存储服务的 URL，用于访问存储桶。                                                     | `http://tos-s3-ap-southeast-1.ibytepluses.com`                                           |
| `AK`              | 访问对象存储服务的 Access Key          | `your-access-key`                                                                        |
| `SK`              | 访问对象存储服务的 Secret Key          | `your-secret-key`                                                                        |

## 使用步骤

1. **准备 `.env` 文件**  
   在项目根目录下创建 `.env` 文件，内容如下：
   ```dotenv
   NAMESPACE=default
   IMAGE=lmsysorg/sglang:latest
   GPU=1
   BUCKET=demo-bucket
   MODEL_PATH_TOS=/models/
   URL=http://tos-s3-ap-southeast-1.ibytepluses.com
   AK=your-access-key
   SK=your-secret-key
   ```

2. **运行脚本**  
   执行以下命令生成最终的 YAML 文件：
   ```bash
   cd volc-examples
   bash scripts/generate_yaml.sh
   ```

3. **查看生成的 YAML 文件**  
   脚本会在当前目录下生成 `final-all-in-one.yaml` 文件，内容为替换后的 Kubernetes 配置。

4. **应用 YAML 文件**  
   使用 `kubectl` 命令将生成的 YAML 文件应用到 Kubernetes 集群：
   ```bash
   kubectl apply -f final-all-in-one.yaml
   ```

## 注意事项

- 确保 `.env` 文件中的变量值正确无误。
- 如果需要自定义文件路径，可以通过设置以下环境变量覆盖默认值：
  - `ENV_FILE`：`.env` 文件路径（默认 `./.env`）。
  - `TEMPLATE_FILE`：模板文件路径（默认 `./all-in-one.yaml`）。
  - `OUTPUT_FILE`：生成的 YAML 文件路径（默认 `./final-all-in-one.yaml`）。
  示例：
  ```bash
  ENV_FILE=custom.env TEMPLATE_FILE=custom-template.yaml OUTPUT_FILE=custom-output.yaml bash scripts/generate_yaml.sh
  ```

- `AK` 和 `SK` 会自动进行 Base64 编码，请确保在 `.env` 文件中提供未编码的值。

运行脚本后，生成的 `final-all-in-one.yaml` 文件将替换所有占位符为实际值。