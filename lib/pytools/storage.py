#!/usr/bin/python
import logging
import subprocess
import json

logging.basicConfig()
_LOGGER = logging.getLogger(__name__)


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


def _get_devices_json():
    cmd = ['lsblk', '-o', 'NAME,SERIAL,TYPE,MODEL,ROTA,UUID,SIZE', '-p', '-J', '-b']
    return _must_succeed(cmd).strip()


def _from_os_json_to_internal_json(os_json_data):
    result_json_data = {}
    block_devices_data = os_json_data['blockdevices']
    result_json_data['disks'] = []
    for element_data in block_devices_data:
        if element_data['type'] == 'disk':
            result_element_json_data = dict()
            result_element_json_data['path'] = element_data['name']
            result_element_json_data['serial'] = element_data['serial']
            media_type = 'SSD' if int(element_data['rota']) is 0 else 'HDD'
            result_element_json_data['mediaType'] = media_type
            result_element_json_data['model'] = element_data['model']
            result_element_json_data['totalCapacityMB'] = int(int(element_data['size']) / 1000000)
            result_json_data['disks'].append(result_element_json_data)
    return result_json_data


def get_storage_list():
    os_json = _get_devices_json()
    os_json_data = json.loads(os_json)
    result_json_data = _from_os_json_to_internal_json(os_json_data)
    return result_json_data


if __name__ == '__main__':  # EXEMPT_FROM_CODE_COVERAGE
    print(get_storage_list())
