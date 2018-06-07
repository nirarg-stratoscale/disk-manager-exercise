import logging
import subprocess

logging.basicConfig()
_LOGGER = logging.getLogger(__name__)

from subsystem import base

# Module will be introduced via PYTHONPATH change
import disk_manager_exercise_client  # pylint: disable=import-error


class TestDisks(base.BaseTest):
    """Endpoints sanity tests."""

    def setUp(self):
        """Wait for the required services to startup."""
        super(TestDisks, self).setUp()

        host_name = _get_hostname().decode('UTF-8')
        self.disk1 = disk_manager_exercise_client.Disk(hostname=host_name,
                                                  id='bc491ea3-02b7-3dea-b0d0-9ac8e7f599b5',
                                                  media_type='SSD',
                                                  model='SAMSUNG TEST1 MZ7TN512',
                                                  path='/dev/sda1',
                                                  serial='S35NNY0HA050941',
                                                  total_capacity_mb=512110111)
        self.disk2 = disk_manager_exercise_client.Disk(hostname=host_name,
                                                  id='9ea18b8b-2cf2-34a2-8c41-d70ecb6efafc',
                                                  media_type='HDD',
                                                  model='SAMSUNG TEST2 MZ7TN512',
                                                  path='/dev/sda2',
                                                  serial='S35NNY0HA050942',
                                                  total_capacity_mb=512110222)
        self.disk3 = disk_manager_exercise_client.Disk(hostname=host_name,
                                                  id='77b84802-5c07-3f12-beb8-5f5828197f98',
                                                  media_type='SSD',
                                                  model='SAMSUNG TEST3 MZ7TN512',
                                                  path='/dev/sda3',
                                                  serial='S35NNY0HA050943',
                                                  total_capacity_mb=512110333)

    def test_disks_list(self):
        """Validate service api."""
        logging.info('Validating list_disks api response')

        result = self.client.list_disks()
        expected = [self.disk1, self.disk2, self.disk3]

        self.assertEquals(result, expected)

    def test_disk_by_id_ok(self):
        """Validate service api."""
        logging.info('Validating disk_by_id api response')
        result = self.client.disk_by_id(disk_id='bc491ea3-02b7-3dea-b0d0-9ac8e7f599b5')

        self.assertEquals(result, self.disk1)

    def test_disk_by_id_wrong_id(self):
        logging.info('Validating disk_by_id with wrong id api response')
        with self.assertRaises(disk_manager_exercise_client.rest.ApiException) as context:
            self.client.disk_by_id(disk_id='wrong-id-test')

        self.assertEquals(context.exception.status, 400)


def _check_output_log_if_needed(*args, **kwargs):  # EXEMPT_FROM_CODE_COVERAGE
    command = kwargs["args"] if "args" in kwargs else args[0]
    _LOGGER.info(str(command))
    output = subprocess.check_output(*args,
                                     stderr=subprocess.STDOUT,
                                     close_fds=True,
                                     **kwargs)
    _LOGGER.debug("Command '%(command)s' finished successfully. Output was:\n%(output)s",
                  dict(command=command, output=output))
    return output


def _must_succeed(*args, **kwargs):  # EXEMPT_FROM_CODE_COVERAGE
    try:
        return _check_output_log_if_needed(*args, **kwargs)
    except subprocess.CalledProcessError as error1:
        _LOGGER.exception("Output was:\n%(output)s\nReturn code:%(returncode)d",
                          dict(output=error1.output, returncode=error1.returncode))
        raise


def _get_hostname():
    cmd = ['hostname']
    return _must_succeed(cmd).strip()