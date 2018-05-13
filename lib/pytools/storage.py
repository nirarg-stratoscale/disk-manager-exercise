#!/usr/bin/python
import logging
import subprocess

logging.basicConfig()
_LOGGER = logging.getLogger(__name__)

HOST_PROC_DIR = "hostproc"
MOUNT_COMMAND = "--mount=/" + HOST_PROC_DIR + "/1/ns/mnt"


def get_storage_list():
    os_output = _get_devices_json()
    result_json_data = _from_os_output_to_internal_json(os_output)
    # os_json_data = json.loads(os_list)
    # result_json_data = _from_os_json_to_internal_json(os_json_data)
    return result_json_data


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
    cmd = ['nsenter', MOUNT_COMMAND, 'lsblk', '-o', 'NAME,SERIAL,TYPE,ROTA,SIZE,MODEL', '-p', '-blnd', '-e', '4,7,10,11,43']
    return _must_succeed(cmd).strip()


def _from_byte_to_string(byte_input):
    result = byte_input.decode("utf-8")
    return result


def _from_os_output_element_to_internal_json(output_element):
    result_element_json_data = dict()
    result_element_json_data['path'] = _from_byte_to_string(output_element[0])
    result_element_json_data['serial'] = _from_byte_to_string(output_element[1])
    media_type = 'SSD' if int(_from_byte_to_string(output_element[3])) is 0 else 'HDD'
    result_element_json_data['mediaType'] = media_type
    model = str(_from_byte_to_string(output_element[5]))
    if len(output_element) > 6:
        for i in range(6, len(output_element)):
            model = model + " " + _from_byte_to_string(output_element[i])
    result_element_json_data['model'] = model
    result_element_json_data['totalCapacityMB'] = int(int(output_element[4]) / 1000000)
    return result_element_json_data


def _from_os_output_to_internal_json(os_output):
    result_json_data = dict()
    result_json_data['disks'] = []
    for line in os_output.splitlines():
        output_element = line.split()
        element_internal_json = _from_os_output_element_to_internal_json(output_element)
        result_json_data['disks'].append(element_internal_json)
    return result_json_data


if __name__ == '__main__':  # EXEMPT_FROM_CODE_COVERAGE
    print(get_storage_list())
