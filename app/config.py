import os


class Settings:
    MAIN_SERVICE_URL = os.getenv("MAIN_SERVICE_URL", "http://10.111.255.45:8082")
    SECRET_TOKEN = os.getenv("SECRET_TOKEN", "glotto88")

    @property
    def callback_url(self) -> str:
        return f"{self.MAIN_SERVICE_URL}/api/internal/lang-calculation"


settings = Settings()
