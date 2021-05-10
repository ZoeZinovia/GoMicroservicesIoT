#include "readDht11.h"
#include <wiringPi.h>
#include <stdio.h>
#include <stdlib.h>
#include <string.h>

// Pi dht11 variables
#define MAXTIMINGS	85
#define DHTPIN		7

int dht11_dat[5] = { 0, 0, 0, 0, 0 };

// Reading of the dht11 is rather complex in C/C++. See this site that explains how readings are made: http://www.uugear.com/portfolio/dht11-humidity-temperature-sensor-module/
int* read_dht11_dat()
{
    u_int8_t laststate	= HIGH;
    u_int8_t counter		= 0;
    u_int8_t j		= 0, i;

    dht11_dat[0] = dht11_dat[1] = dht11_dat[2] = dht11_dat[3] = dht11_dat[4] = 0;

    // pull pin down for 18 milliseconds. This is called “Start Signal” and it is to ensure DHT11 has detected the signal from MCU.
    pinMode( DHTPIN, OUTPUT );
    digitalWrite( DHTPIN, LOW );
    delay( 18 );
    // Then MCU will pull up DATA pin for 40us to wait for DHT11’s response.
    digitalWrite( DHTPIN, HIGH );
    delayMicroseconds( 40 );
    // Prepare to read the pin
    pinMode( DHTPIN, INPUT );

    // Detect change and read data
    for ( i = 0; i < MAXTIMINGS; i++ )
    {
        counter = 0;
        while ( digitalRead( DHTPIN ) == laststate )
        {
            counter++;
            delayMicroseconds( 1 );
            if ( counter == 255 )
            {
                break;
            }
        }
        laststate = digitalRead( DHTPIN );

        if ( counter == 255 )
            break;

        // Ignore first 3 transitions
        if ( (i >= 4) && (i % 2 == 0) )
        {
            // Add each bit into the storage bytes
            dht11_dat[j / 8] <<= 1;
            if ( counter > 16 )
                dht11_dat[j / 8] |= 1;
            j++;
        }
    }

    // Check that 40 bits (8bit x 5 ) were read + verify checksum in the last byte
    if ( (j >= 40) && (dht11_dat[4] == ( (dht11_dat[0] + dht11_dat[1] + dht11_dat[2] + dht11_dat[3]) & 0xFF) ) )
    {
        return dht11_dat; // If all ok, return pointer to the data array
    } else  {
        dht11_dat[0] = -1;
        return dht11_dat; //If there was an error, set first array element to -1 as flag to main function
    }
}

