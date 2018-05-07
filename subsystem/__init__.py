"""Monkey patch nose2 log-capture plugin to use the required log formatter.

See https://github.com/nose-devs/nose2/issues/135 for more details.
"""
import logging

from nose2.plugins.logcapture import LogCapture

# Define the required log format for the tests
FORMATTER = logging.Formatter(
    fmt='%(asctime)s.%(msecs)03d %(levelname)-8s%(message)s [%(name)s]',
    datefmt="%YY-%m-%dT%H:%M:%S"
)

start_test = LogCapture.startTest


def patched_start_test(self, event):
    """Override log handler formatter then execute original method."""
    self.handler.setFormatter(FORMATTER)
    return start_test(self, event)


LogCapture.startTest = patched_start_test
