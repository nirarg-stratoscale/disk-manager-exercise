import logging
from unittest import skip

from subsystem import base

# Module will be introduced via PYTHONPATH change
import disk_manager_exercise_client  # pylint: disable=import-error


class TestDisks(base.BaseTest):
    """Endpoints sanity tests."""

    def test_disks_list(self):
        """Validate service api."""
        logging.info('Validating api response')
        import ipdb
        ipdb.set_trace()
        with self.assertRaises(disk_manager_exercise_client.rest.ApiException) as context:
            # TODO: Use client lib to invoke the API methods, e.g
            self.client.list_disks()

        self.assertEquals(context.exception.status, 200)
        # self.assertEquals(context.exception.reason, 'Not Implemented')
