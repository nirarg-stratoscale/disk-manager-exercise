#!/usr/bin/python
import unittest
import mock as mock
import pytools.storage


class TestStorage(unittest.TestCase):

    @mock.patch('subprocess.check_output')
    def test_get_storage_list(self, check_output):
        # Test one result, one SSD disk
        check_output.return_value = '''
                                    {
                                       "blockdevices": [
                                          {"name": "/dev/sda", "serial": "S35NNY0HA05094", "type": "disk", "model": "SAMSUNG MZ7TN512", "rota": "0", "uuid": null, "size": "512110190592"}
                                       ]
                                    }
                                    '''
        expected = {'disks': [{'mediaType': 'SSD', 'path': u'/dev/sda', 'serial': u'S35NNY0HA05094',
                               'totalCapacityMB': 512110, 'model': u'SAMSUNG MZ7TN512'}]}
        self.assertEqual(pytools.storage.get_storage_list(), expected)

        # Test one result, one HDD disk
        check_output.return_value = '''
                                    {
                                       "blockdevices": [
                                          {"name": "/dev/sda", "serial": "S35NNY0HA05094", "type": "disk", "model": "SAMSUNG MZ7TN512", "rota": "1", "uuid": null, "size": "512110190592"}
                                       ]
                                    }
                                    '''
        expected = {'disks': [
            {'mediaType': 'HDD', 'path': u'/dev/sda', 'serial': u'S35NNY0HA05094',
             'totalCapacityMB': 512110, 'model': u'SAMSUNG MZ7TN512'}]}
        self.assertEqual(pytools.storage.get_storage_list(), expected)

        # Test empty result, No disks
        check_output.return_value = '''{"blockdevices": []}'''
        expected = {'disks': []}

        self.assertEqual(pytools.storage.get_storage_list(), expected)

        # Test three results, one disk
        check_output.return_value = '''
                                    {
                                       "blockdevices": [
                                          {"name": "/dev/loop0", "serial": null, "type": "loop", "model": null, "rota": "1", "uuid": null, "size": "98062336"},
                                          {"name": "/dev/loop1", "serial": null, "type": "loop", "model": null, "rota": "1", "uuid": null, "size": "169943040"},
                                          {"name": "/dev/sda", "serial": "S35NNY0HA05094", "type": "disk", "model": "SAMSUNG MZ7TN512", "rota": "0", "uuid": null, "size": "512110190592",
                                             "children": [
                                                {"name": "/dev/sda1", "serial": null, "type": "part", "model": null, "rota": "0", "uuid": "E2BB-791A", "size": "209715200"},
                                                {"name": "/dev/sda2", "serial": null, "type": "part", "model": null, "rota": "0", "uuid": "f667a867-6628-42c0-9fd2-3dc3b115da5e", "size": "1073741824"},
                                                {"name": "/dev/sda3", "serial": null, "type": "part", "model": null, "rota": "0", "uuid": "LavYOD-YHKh-W97t-SWCz-1CIE-3tDP-LCA0lA", "size": "510825332736",
                                                   "children": [
                                                      {"name": "/dev/mapper/fedora-root", "serial": null, "type": "lvm", "model": null, "rota": "0", "uuid": "52417c50-d832-4fb0-9d15-fef48a482ac7", "size": "53687091200"},
                                                      {"name": "/dev/mapper/fedora-swap", "serial": null, "type": "lvm", "model": null, "rota": "0", "uuid": "c9785206-fb7a-4371-a68f-cf6a53220dac", "size": "8422162432"},
                                                      {"name": "/dev/mapper/fedora-home", "serial": null, "type": "lvm", "model": null, "rota": "0", "uuid": "6652c78a-ac48-4c78-84e1-67b5fc6abff7", "size": "448715030528"}
                                                   ]
                                                }
                                             ]
                                          }
                                       ]
                                    }
                                    '''
        expected = {'disks': [
            {'mediaType': 'SSD', 'path': u'/dev/sda', 'serial': u'S35NNY0HA05094',
             'totalCapacityMB': 512110, 'model': u'SAMSUNG MZ7TN512'}]}
        self.assertEqual(pytools.storage.get_storage_list(), expected)

        # Test two results, no disks
        check_output.return_value = '''
                                    {
                                       "blockdevices": [
                                          {"name": "/dev/loop0", "serial": null, "type": "loop", "model": null, "rota": "1", "uuid": null, "size": "98062336"},
                                          {"name": "/dev/loop1", "serial": null, "type": "loop", "model": null, "rota": "1", "uuid": null, "size": "169943040"}
                                       ]
                                    }
                                    '''
        expected = {'disks': []}
        self.assertEqual(pytools.storage.get_storage_list(), expected)

        # Test one result, error
        check_output.return_value = '''
                                    {
                                       "blockdevices": [
                                          {"name": "/dev/sda", "type": "disk", "model": "SAMSUNG MZ7TN512", "rota": "1", "uuid": null, "size": "512110190592"}
                                       ]
                                    }
                                    '''
        self.assertRaises(KeyError, pytools.storage.get_storage_list)
