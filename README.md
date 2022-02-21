<p align="center">
  <a href="https://github.com/fwiedmann/icof">
    <img src="icof.jpeg" width=100 height=100>
  </a>

<h3 align="center">icof</h3>

  <p align="center">
     in case of fire
  </p>
</p>

![love](https://img.shields.io/badge/made%20with-%E2%9D%A4%EF%B8%8F-lightgrey)


Due to the ongoing global pandemic I started to work mostly from home.
Beside my job as a software engineer I am also a volunteer firefighter in my local community.
This means that I'm on call for emergencies 24/7, also during my daily work as a SE.
When I am in the office and get a call, I can speak to my colleagues in seconds that I have to go (Or they just see me running to my car).
In booth cases they are informed that I am afk and can probably re-schedule or take over meetings for me.
When working from home I want to send an absence notification to my colleagues.
Chatting my 2-4 colleagues can be time-consuming.
Well, I am an engineer. I can automate this in combination with a big red button. I like big red buttons.



## gocrazy installation

First you need to plug-in the Pi SD card into your PC. If you have already installed an OS like raspbian on this SD card, I recommend flashing it.
With the `fdisk` utility you can delete all existing partitions.

Now you can follow the go-krazy installation guide: https://github.com/gokrazy/gokrazy#installation

You can update your gokrazy installation via the web endpoint:
```bash
gokr-packer -update yes -hostname 192.168.178.49 -serial_console=disabled ./cmd/icof-cmd
```

### Configuration

To run icof on gokrazy you have to create a startup configuration and place it in the following directory on the permanent data partition of your gocrazy installation: `/perm/icof/start-config.json`.

Configuration:
```json
{
  "email_config": {
    "host": "",
    "port": 0,
    "username": "",
    "password": "",
    "from_email_address": ""
  },
  "email_receiver_config": {
    "alert_subject": "",
    "resolve_subject": "",
    "receivers": [
      {
        "name": "",
        "alert_template_message": "",
        "resolve_template_message": "",
        "addresses": [
          {
            "email": "",
            "name": "",
            "surname": ""
          }
        ]
      }
    ]
  }
}
```

## Circuit 

The circuit uses a pull down resistor. This will pull down the input signal at the GPIO pin to LOW if the switch is pressed and the power circuit is interrupted. 

- Emergency switch is wired in NC (normally closed).
- A 10kΩ resistor is used
- 3,3V power pin is enough

![circuit](./icof_circuit.png)

## Shopping List

- A raspberry pi 3 B+ or 4 (choose whats curently available) https://www.amazon.de/s?k=raspberry+pi&__mk_de_DE=%C3%85M%C3%85%C5%BD%C3%95%C3%91&crid=2J4U87EF8YB01&sprefix=raspberry+pi+%2Caps%2C126&ref=nb_sb_noss
- Jumper cables + Breadboard https://www.amazon.de/AZDelivery-%E2%AD%90%E2%AD%90%E2%AD%90%E2%AD%90%E2%AD%90-Jumper-Breadboard-Arduino/dp/B078JGQKWP/ref=sr_1_4?__mk_de_DE=%C3%85M%C3%85%C5%BD%C3%95%C3%91&crid=1QJZY9KRGO4MI&keywords=azdelivery+set&qid=1645385067&sprefix=azdelivery+set%2Caps%2C87&sr=8-4
- Resistors https://www.amazon.de/AZDelivery-Widerst%C3%A4nde-Resistor-Widerstand-Sortiment/dp/B07Q87JZ9G?ref_=ast_sto_dp&th=1&psc=1