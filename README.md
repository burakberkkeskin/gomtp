## 📬 GOMTP

Gompt is a cli tool to test smtp settings easily.

## Install

sudo curl -L -o /usr/local/bin/gomtp "https://github.com/safderun/gomtp/releases/latest/download/gomtp-$(uname -s)-$(uname -m)" && \
sudo chmod +x /usr/local/bin/gomtp

## Usage

- Create a `gomtp.yml` file anywhre you want.
- Take the template from the `gomtp.yml`
- There is 4 templates for `mailhog`, `gmail`, `yandex` and `brevo`
- `subject` and `body` is optional.
- In the same directory with your configured `gomtp.yml`, run `gomtp` with no argument.

```bash
❯ gompt
Email sent successfully!
```

- If your configuration is valid, you will see the "Email sent successfully!" message.

## Sample SMTP For Testing

- To test the `gomtp` quickly, you can run the `mailhog` from `docker-compose.yml`

```bash
docker compose up -d
```

- Configure `gomtp.yml`
- Open the `mailhog` web ui from http://127.0.0.1:8025
- Run the `gomtp`.
