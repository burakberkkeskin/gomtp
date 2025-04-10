## 📬 GOMTP

Gomtp is a cli tool to test smtp settings easily.

## Install

### Install From Binary (Recommended)

You can install the `gomtp` to Linux or macOS with these commands:

```bash
sudo curl -L -o /usr/local/bin/gomtp "https://github.com/burakberkkeskin/gomtp/releases/latest/download/gomtp-$(uname -s)-$(uname -m)" && \
sudo chmod +x /usr/local/bin/gomtp
```

### Build Locally

You can build the `gomtp` locally, on your own machine.

```bash
version=$(git describe --tags --abbrev=0) && \
commitId=$(git --no-pager log -1 --oneline | awk '{print $1}') && \
go build -ldflags "-X gomtp/cmd.version=$version -X gomtp/cmd.commitId=$commitId" -o gomtp -v .
```

## Usage

- Create a `gomtp.yaml` file anywhre you want.
- Take the template from the `gomtp.yaml`
- There is 4 templates for `mailhog`, `gmail`, `yandex` and `brevo`
- `subject` and `body` is optional.
- In the same directory with your configured `gomtp.yaml`, run `gomtp` with no argument.

```bash
❯ gomtp
Email sent successfully!
```

- If your configuration is valid, you will see the "Email sent successfully!" message.

## Custom Gomtp Yaml Path

- You can name the `gomtp.yaml` as you wish while creating the configuration.

- If you change the default configuration file name, you can pass the path of the file to the `gomtp`.

```bash
gomtp --file test.yaml
```

or

```bash
gomtp -f test.yaml
```

## Sample SMTP For Testing

To test the `gomtp` quickly, you can run the `mailpit` from `docker-compose.yml`

- Install the `gomtp` by checking the [Install](#install) section.

- Create a separate directory for config files

```bash
mkdir ~/gomtp 
```

- Change directory

```bash
cd ~/gomtp/ 
```

- Set the gomtp binary permissions

```bash
sudo chmod +x /usr/local/bin/gomtp
```

- Download the sample docker-compose.yaml

```bash
curl -LO https://raw.githubusercontent.com/burakberkkeskin/gomtp/refs/heads/master/docker-compose.yaml
```

- Download the sample gomtp.yaml configuration.
  - The default `gomtp.yaml` file already has been configured for the `mailpit`.

```bash
curl -LO https://raw.githubusercontent.com/burakberkkeskin/gomtp/refs/heads/master/gomtp.yaml
```

- Start the mailpit.

```bash
docker compose up -d
```

- Test the gomtp

```bash
gomtp
```

```output
Email sent successfully!
```

- Open the `mailpit` web ui from http://127.0.0.1:8025 and see the sample email.

## Configure Once, Use For Anything

If you want to use `gomtp` to send emails, you can configure a yaml and use it as base. For example, follow the use-case below: 

- Create the template file for `gomtp.yaml` on your home directory.
- Configure the `gomtp.yaml` for your Gmail account in `~/gomtp.yaml` path.
```bash
vim ~/gomtp.yaml
```
- Configure the `username`, `password`, `from`. Optionally delete the `to`, `subject` and `body` fields.
- Now use `gomtp` use send email to any email address.
- With `--body` flag:
```bash
gomtp -f ~/gomtp.yaml --to yourTargetEmailAddress@gmail.com --subject "Test Email From gomtp" --body "test from atlantic server"
```
- With `--body-file` flag:
```bash
gomtp -f ~/gomtp.yaml --to yourTargetEmailAddress@gmail.com --subject "Test Email From gomtp" --body-file "~/email-body.log"
```
- With piping a command output:
```bash
echo "This is an example command or bash script output" | gomtp -f ~/gomtp.yaml --to yourTargetEmailAddress@gmail.com --subject "Test Email From gomtp"
```

## Release a version

- Define a version.

```bash
export gomtpVersion=v1.4.0
```

- You should create a release branch from the master

```bash
git checkout master && git pull && \
git checkout -b release/${gomtpVersion}
```

- Tag the commit

```bash
git tag --sign ${gomtpVersion} -m "Added verifyCertificate and example commands."
```

- Push the release branch and tags

```bash
git push && git push --tags
```

## Run Tests

- Before run the tests, run the docker compose file.

```bash
docker compose up -d
```

- You can run e2e tests to ensure application stability.

```bash
go test -p 1 ./cmd/
```

- Check the test coverage:

```bash
go test -p 1 ./cmd/ -coverprofile=coverage.out
```

- You can see covered lines with html report:

```bash
go test -p 1 ./cmd/ -coverprofile=coverage.out  -test.coverprofile ./c.out && \
go tool cover -html=c.out
```
