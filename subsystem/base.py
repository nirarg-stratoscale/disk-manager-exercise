from docker_test_tools import base_test

# Module will be introduced via PYTHONPATH change
import disk_manager_exercise_client  # pylint: disable=import-error


class BaseTest(base_test.BaseDockerTest):
    """Basic subsystem test."""
    CHECKS_TIMEOUT = 60  # Used by docker-test-tools for health checks that services are up

    ENDPOINT = 'http://localhost:80/api/v2'
    DEFAULT_USER_ID = 'user_id'
    DEFAULT_PROJECT_ID = 'project_id'
    DEFAULT_DOMAIN_ID = 'domain_id'
    DEFAULT_ROLES = 'admin'

    def setUp(self):
        """Wait for the required services to startup."""
        super(BaseTest, self).setUp()

        api_client_config = disk_manager_exercise_client.Configuration()
        api_client_config.host = self.ENDPOINT

        api_client=disk_manager_exercise_client.ApiClient(configuration=api_client_config)
        api_client.default_headers['X-Auth-Project-Id'] = self.DEFAULT_PROJECT_ID
        api_client.default_headers['X-Auth-User-Id'] = self.DEFAULT_USER_ID
        api_client.default_headers['X-Auth-User-Domain-Id'] = self.DEFAULT_DOMAIN_ID
        api_client.default_headers['X-Auth-Roles'] = self.DEFAULT_ROLES

        self.client = disk_manager_exercise_client.api.DiskApi(api_client=api_client)
