import RPi.GPIO as GPIO
import argparse

parser = argparse.ArgumentParser(description='Control a GPIO pin.')
parser.add_argument('-n', default='18', type=int, help='The number of the GPIO pin.')
parser.add_argument('-a', default='on', help='The action to perform. ("on", "off" or "toggle")')
args = parser.parse_args()

GPIO.setmode(GPIO.BCM)
GPIO.setwarnings(False)
GPIO.setup(args.n,GPIO.OUT)
if (args.a == 'on'):
    GPIO.output(args.n, GPIO.HIGH)    
elif (args.a == 'toggle'):
    if GPIO.input(args.n):
        GPIO.output(args.n, GPIO.LOW)
    else:
        GPIO.output(args.n, GPIO.HIGH)
else:
    GPIO.output(args.n, GPIO.LOW)
