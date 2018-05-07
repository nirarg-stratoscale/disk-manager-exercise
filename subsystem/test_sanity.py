import logging
from unittest import skip

from subsystem import base

# Module will be introduced via PYTHONPATH change
import disk_manager_exercise_client  # pylint: disable=import-error


class TestSanity(base.BaseTest):
    """Endpoints sanity tests."""

    @skip("Not implemented")
    def test_sanity(self):
        """Validate service api."""
        logging.info('Validating api response')
        with self.assertRaises(disk_manager_exercise_client.rest.ApiException) as context:
            # TODO: Use client lib to invoke the API methods, e.g
            # self.client.operation_id()
            pass

        self.assertEquals(context.exception.status, 501)
        self.assertEquals(context.exception.reason, 'Not Implemented')