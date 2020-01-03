from gpiozero import LightSensor
from time import sleep

pulseCount = 0

def lightPulse():
    global pulseCount
    pulseCount = pulseCount + 1
    print(pulseCount)

ldr = LightSensor(19,queue_len=1)
ldr.when_light = lightPulse
ldr.threshold = 0.1

while True:
    sleep(1)