from gpiozero import LED
from time import sleep

led = LED(26)

for x in range(20):
    led.on()
    sleep(0.1)
    led.off()
    print("1")
    sleep(0.3)