from docker_test_tools import base_test

# Module will be introduced via PYTHONPATH change
import disk_manager_exercise_client  # pylint: disable=import-error


class BaseTest(base_test.BaseDockerTest):
    """Basic subsystem test."""
    CHECKS_TIMEOUT = 60  # Used by docker-test-tools for health checks that services are up

    ENDPOINT = 'http://localhost:80/api/v2'

    def setUp(self):
        """Wait for the required services to startup."""
        super(BaseTest, self).setUp()

        api_client_config = disk_manager_exercise_client.Configuration()
        api_client_config.host = self.ENDPOINT

        self.client = disk_manager_exercise_client.api.Disk-manager-exerciseApi(
            api_client=disk_manager_exercise_client.ApiClient(
                configuration=api_client_config
            )
        )
