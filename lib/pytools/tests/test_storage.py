#!/usr/bin/python
import unittest
import mock as mock
import pytools.storage


class TestStorage(unittest.TestCase):

    @mock.patch('subprocess.check_output')
    def test_storage_basic(self, check_output):
        # Test one result, one SSD disk
        check_output.return_value = b'/dev/sda S35NNY0HA05094 disk    0 512110190592 SAMSUNG MZ7TN512'
        expected = {'disks': [{'mediaType': 'SSD', 'path': u'/dev/sda', 'serial': u'S35NNY0HA05094',
                               'totalCapacityMB': 512110, 'model': u'SAMSUNG MZ7TN512'}]}
        self.assertEqual(pytools.storage.get_storage_list(), expected)


    @mock.patch('subprocess.check_output')
    def test_torage_hdd(self, check_output):
        # Test one result, one SSD disk
        check_output.return_value = b'/dev/sda S35NNY0HA05094 disk    1 512110190592 SAMSUNG MZ7TN512'
        expected = {'disks': [{'mediaType': 'HDD', 'path': u'/dev/sda', 'serial': u'S35NNY0HA05094',
                               'totalCapacityMB': 512110, 'model': u'SAMSUNG MZ7TN512'}]}
        self.assertEqual(pytools.storage.get_storage_list(), expect


    @mock.patch('subprocess.check_output')
    def test_storage_model_one_word(self, check_output):
        # Test one result, one SSD disk
        check_output.return_value = b'/dev/sda S35NNY0HA05094 disk    1 512110190592 SAMSUNG'
        expected = {'disks': [{'mediaType': 'HDD', 'path': u'/dev/sda', 'serial': u'S35NNY0HA05094',
                               'totalCapacityMB': 512110, 'model': u'SAMSUNG'}]}
        self.assertEqual(pytools.storage.get_storage_list(), expected)


    @mock.patch('subprocess.check_output')
    def test_storage_model_three_words(self, check_output):
        # Test one result, one SSD disk
        check_output.return_value = b'/dev/sda S35NNY0HA05094 disk    1 512110190592 SAMSUNG MZ7TN512 ASDF'
        expected = {'disks': [{'mediaType': 'HDD', 'path': u'/dev/sda', 'serial': u'S35NNY0HA05094',
                               'totalCapacityMB': 512110, 'model': u'SAMSUNG MZ7TN512 ASDF'}]}
        self.assertEqual(pytools.storage.get_storage_list(), expected)


    @mock.patch('subprocess.check_output')
    def test_storage_model_special(self, check_output):
        # Test one result, one SSD disk
        check_output.return_value = b'/dev/sda S35NNY0HA05094 disk    1 512110190592 !@#$%'
        expected = {'disks': [{'mediaType': 'HDD', 'path': u'/dev/sda', 'serial': u'S35NNY0HA05094',
                               'totalCapacityMB': 512110, 'model': u'!@#$%'}]}
        self.assertEqual(pytools.storage.get_storage_list(), expected)


    @mock.patch('subprocess.check_output')
    def test_storage_model_empty(self, check_output):
        # Test one result, one SSD disk
        check_output.return_value = b''
        expected = {'disks': []}
        self.assertEqual(pytools.storage.get_storage_list(), expected)


    @mock.patch('subprocess.check_output')
    def test_storage_model_two(self, check_output):
        # Test one result, one SSD disk
        check_output.return_value = b'/dev/sda1 S35NNY0HA0509411 disk 1 10000000 ASDFGH\n' \
                                    b'/dev/sda2 S35NNY0HA0509422 disk    0 512110190592 QWERTY'
        expected = {'disks': [{'mediaType': 'HDD', 'path': u'/dev/sda1',
                               'serial': u'S35NNY0HA0509411', 'totalCapacityMB': 10, 'model': u'ASDFGH'},
                              {'mediaType': 'SSD', 'path': u'/dev/sda2', 'serial': u'S35NNY0HA0509422',
                               'totalCapacityMB': 512110, 'model': u'QWERTY'}]}
        self.assertEqual(pytools.storage.get_storage_list(), expected)
