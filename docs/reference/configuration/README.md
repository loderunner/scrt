# Configuration

## Global

### Password

- Type: `string`
- YAML: `password`
- Environment variable: `SCRT_PASSWORD`

The password to the store. The setting will be used to derive a key, to decrypt and encrypt the data in the store.

### Storage type

- Type: `string`, `"local" | "s3" | "git"`
- YAML: N/A
- Environment variable: `SCRT_STORAGE`

The storage backend to use for the store.

Storage type (`storage`) can be ignored in the YAML configuration file. scrt will read the configuration under the key for the storage type (e.g. `local:`). _Defining configurations for multiple storage types in a single file will result in undefined behavior._

### Verbosity

- Type: `boolean`
- Default: `false`
- YAML: `verbose`
- Environment variables: `SCRT_VERBOSE`

## Local storage

### Path

- Type: `string`
- YAML: `local` > `path`
- Environment variable: `SCRT_LOCAL_PATH`

The path to the store file on the local computer.

## S3 storage

### Bucket name

- Type: `string`
- YAML: `s3` > `bucket-name`
- Environment variable: `SCRT_S3_BUCKET_NAME`

The name of the bucket where the store object is located.

### Object key

- Type: `string`
- YAML: `s3` > `key`
- Environment variable: `SCRT_S3_KEY`

The path of the store object in the bucket.

### Endpoint URL

- Type: `string`
- Default: `https://s3.<region>.amazonaws.com`
- YAML: `s3` > `endpoint-url`
- Environment variable: `SCRT_S3_ENDPOINT_URL`

Override the default S3 URL.

### Region

- Type: `string`
- YAML: `s3` > `region`
- Environment variable: `SCRT_S3_REGION`

The region of the S3 storage.

## Git storage

### URL

- Type: `string`
- YAML: `git` > `url`
- Environment variable: `SCRT_GIT_URL`

The URL of the git repository.

### Path

- Type: `string`
- YAML: `git` > `path`
- Environment variables: `SCRT_GIT_PATH`

The path of the store file in the repository.

### Branch

- Type: `string`
- YAML: `git` > `branch`
- Environment variables: `SCRT_GIT_BRANCH`

The name of the branch to checkout, commit and push to on updates. Uses the default branch when missing.

### Commit or tag

- Type: `string`
- YAML: `git` > `checkout`
- Environment variables: `SCRT_GIT_CHECKOUT`

A tree-ish revision to checkout, e.g. commit or tag

### Commit message

- Type: `string`
- Default: `"update secrets"`
- YAML: `git` > `message`
- Environment variables: `SCRT_GIT_MESSAGE`

The commit message used when updating the store.
