/******************************************************************************
  /// @file       dcdriver.cpp
  /// @date       22.11.2018
  /// @copyright  Copyright (C) beyond protocol inc. - All Rights Reserved
  ///             You may use, distribute and modify this code under the
  ///             terms and conditions defined in the file 'LICENSE.txt',
  ///             which is part of this source code package.
  /// @brief    Raspberry Pi DS28C36 I2C interface driver
  /// @details
  ///
  /// The board was connected as follows:
  /// (Raspberry Pi)(DS28C36)
  /// GND  -> GND
  /// 3.3V -> Vcc
  /// SCL  -> SCL
  /// SDA  -> SDA

  To manualy compile dcdriver as shared (.so) and static (.a) object

  1. Compile source
  gcc -c -fPIC dcdriver.cpp -o libdcdriver.o -O3 -s

  2. Build shared object
  gcc -shared libdcdriver.o -o libdcdriver.so -O3 -s

  3. Build static object
  ar rcs libdcdriver.a libdcdriver.o
******************************************************************************/

/***
  DS28C36 Memory Setup
  --------------------
   Page 0:  not used, filled with 00h, no protection <BR>
   Page 1:  not used, filled with 00h, no protection <BR>
   Page 2:  not used, filled with 00h, no protection <BR>
   Page 3:  not used, filled with 00h, no protection <BR>
   Page 4:  not used, filled with 00h, no protection <BR>
   Page 5:  not used, filled with 00h, no protection <BR>
   Page 6:  not used, filled with 00h, no protection <BR>
   Page 7:  not used, filled with 00h, no protection <BR>
   Page 8:  not used, filled with 00h, no protection <BR>
   Page 9:  not used, filled with 00h, no protection <BR>
   Page 10: not used, filled with 00h, no protection <BR>
   Page 11: not used, filled with 00h, no protection <BR>
   Page 12: not used, filled with 00h, no protection <BR>
   Page 13: not used, filled with 00h, no protection <BR>
   Page 14: ECDH Certificate r (WP) Slave Pub Key signed w/ Verify Authority Private key <BR>
   Page 15: ECDH Certificate s (WP) Slave Pub Key signed w/ Verify Authority Private key <BR>
   Page 16: Public Key AX, Slave ECC key set, Device Generated (WP) <BR>
   Page 17: Public Key AY, Slave ECC key set, Device Generated (WP) <BR>
   Page 18: Public Key BX, not used, random data (WP) <BR>
   Page 19: Public Key BY, not used, random data (WP) <BR>
   Page 20: Public Key CX, Authority Public Key X (WP) <BR>
   Page 21: Public Key CY, Authority Public Key Y (WP) <BR>
   Page 22: Private Key A, Slave ECC key set, Device Generated (RP,WP) <BR>
   Page 23: Private Key B, not used, random data (RP,WP) <BR>
   Page 24: Private Key C, not used, random data (RP,WP) <BR>
   Page 25: Secret A, not used, random data (RP,WP) <BR>
   Page 26: Secret B, not used, random data (RP,WP) <BR>
   Page 27: Decrement Counter, not set <BR>
   Page 28: ROM Functions, no protection <BR>
   Page 29: GPIO Control, no protection <BR>
   Page 30: Public Key SX, not used <BR>
   Page 31: Public Key SY, not used <BR>

   Certificate Format
   ------------------
   (32-bytes) Public Key X
   (32-bytes) Public Key Y
   (16-bytes) Customization cert field (currently zeros)
   (8-bytes) ROMID
   (2-bytes) MANID
   = Total (90 bytes) Signed by Verify Authority Private Key
*/


#include <stdio.h>
#include <errno.h>
#include <unistd.h>
#include <fcntl.h>    /* For O_RDWR */
#include <sys/ioctl.h>
#include <string.h>

#include "dcdriver.h"

#define I2C_DC_IC_ADR 0x1b //27
#define HEX_SUB_LEN 2 //"0F " or "0F"
#define I2C_SMBUS_READ  1
#define I2C_SMBUS_WRITE 0
#define I2C_SLAVE 0x0703

// DS28C36 commands
#define CMD_WRITE_MEM            0x96
#define CMD_READ_MEM             0x69
#define CMD_WRITE_BUF            0x87
#define CMD_READ_BUF             0x5A
#define CMD_READ_PAGE_PROT       0xAA
#define CMD_SET_PAGE_PROT        0xC3
#define CMD_DECREMENT_CNT        0xC9
#define CMD_READ_RNG             0xD2
#define CMD_ENC_READ_MEM         0x4B
#define CMD_COMP_READ_AUTH       0xA5
#define CMD_AUTH_SHA2_WRITE      0x99
#define CMD_COMP_LOCK_SHA2       0x3C
#define CMD_GEN_ECDSA_KEY        0xCB
#define CMD_COMPUTE_MULI_HASH    0x33
#define CMD_VERIFY_ECDSA_SIG     0x59
#define CMD_AUTH_ECDSA_PUB_KEY   0xA8
#define CMD_AUTH_ECDSA_WRITE     0x89

// Result bytes
#define RESULT_SUCCESS                0xAA //170
#define RESULT_FAIL_PROTECTION        0x55
#define RESULT_FAIL_PARAMETETER       0x77
#define RESULT_FAIL_INVALID_SEQUENCE  0x33
#define RESULT_FAIL_VERIFY            0x00
#define RESULT_FAIL_ECDSA             0x22
#define RESULT_FAIL_COMMUNICATION     0x11

// Special Purpose pages
#define PG_SLAVE_ECDH_CERTIFICATE_R      14
#define PG_SLAVE_ECDH_CERTIFICATE_S      15
#define PG_PUB_KEY_AX            16
#define PG_PUB_KEY_AY            17
#define PG_PUB_KEY_BX            18
#define PG_PUB_KEY_BY            19
#define PG_PUB_KEY_CX            20
#define PG_PUB_KEY_CY            21

#define PG_PRIV_KEY_A            22
#define PG_PRIV_KEY_B            23
#define PG_PRIV_KEY_C            24

#define PG_SECRET_A              25
#define PG_SECRET_B              26

#define PG_DECREMENT_CNT         27
#define PG_ROM_OPTIONS           28
#define PG_GPIO_CONTROL          29

#define PG_PUB_KEY_SX            30
#define PG_PUB_KEY_SY            31

// Offset into GPIO page
#define OFFSET_GPIO_P0C          0x00
#define OFFSET_GPIO_P1C          0x01
#define OFFSET_GPIO_P0L          0x02
#define OFFSET_GPIO_P1L          0x03

// Offset into ROM options page
#define OFFSET_ROM_RBD           0x00
#define OFFEST_ROM_ANON          0x01
#define OFFSET_ROM_MANID         0x22
#define OFFSET_ROM_ROM           0x24

// Protection bit fields
#define PROT_RP                  0x01  // Read Protection  (KEY_PAGES all have RP set by default)
#define PROT_WP                  0x02  // Write Protection
#define PROT_EM                  0x04  // EPROM Emulation Mode (not applicable to KEY_PAGES)
#define PROT_APH                 0x08  // Authentication Write Protection HMAC (not applicable to KEY_PAGES)
#define PROT_EPH                 0x10  // Encryption and Authenticated Write Protection HMAC (not applicable to KEY_PAGES)
#define PROT_AUTH                0x20  // Designated Authority Public Key. Only on AUTH_KEY_PAGE.
#define PROT_ECH                 0x40  // Encrypted read and write using shared key from ECDH
#define PROT_ECW                 0x80  // Authentication Write Protection ECDSA (not applicable to KEY_PAGES)

// AT bit fields for Compute and Read Page Authentication
#define AT_HMAC_SECRETA          0x00  // 000b HMAC using SHA2 Secret A
#define AT_HMAC_SECRETB          0x01  // 001b	HMAC using SHA2 Secret B
#define AT_HMAC_SECRETS          0x02  // 010b	HMAC using SHA2 Secret S
#define AT_ECDSA_KEYA            0x03  // 011b	ECDSA Page Signature using Private Key A
#define AT_ECDSA_KEYB            0x04  // 100b	ECDSA Page Signature using Private Key B
#define AT_ECDSA_KEYC            0x05  // 101b	ECDSA Page Signature using Private Key C (invalid if AUTH protection set)

static unsigned char crc8;
static unsigned char read_buff[256], write_buff[256];
static int i2c_fd;
static int driver_initialized = 0;
static char resTMP[4];
static char buffTMP[256];

char page_name00[18] = "used for sig.  @ ";
char page_name14[18] = "ECDH_Cert_r      ";
char page_name15[18] = "ECDH_Cert_s      ";
char page_name16[18] = "PG_PUB_KEY_AX  @ ";
char page_name17[18] = "PG_PUB_KEY_AY  @ ";
char page_name18[18] = "PG_PUB_KEY_BX    ";
char page_name19[18] = "PG_PUB_KEY_BY    ";
char page_name20[18] = "PG_PUB_KEY_CX    ";
char page_name21[18] = "PG_PUB_KEY_CY    ";
char page_name22[18] = "PG_PRIV_KEY_A* @ ";
char page_name23[18] = "PG_PRIV_KEY_B*   ";
char page_name24[18] = "PG_PRIV_KEY_C*   ";
char page_name25[18] = "PG_SECRET_A*     ";
char page_name26[18] = "PG_SECRET_B*     ";
char page_name27[18] = "PG_DECREMENT_CNT ";
char page_name28[18] = "PG_ROM_OPTIONS @ ";
char page_name29[18] = "PG_GPIO_CONTROL  ";
char page_name30[18] = "PG_PUB_KEY_SX    ";
char page_name31[18] = "PG_PUB_KEY_SY    ";
char page_notusd[18] = "not used         ";

char * printPageName (int page) {
    switch (page) {
        case  0: return page_name00;
        case 14: return page_name14;
        case 15: return page_name15;
        case 16: return page_name16;
        case 17: return page_name17;
        case 18: return page_name18;
        case 19: return page_name19;
        case 20: return page_name20;
        case 21: return page_name21;
        case 22: return page_name22;
        case 23: return page_name23;
        case 24: return page_name24;
        case 25: return page_name25;
        case 26: return page_name26;
        case 27: return page_name27;
        case 28: return page_name28;
        case 29: return page_name29;
        case 30: return page_name30;
        case 31: return page_name31;
        default: return page_notusd;
    }
}

int initDriver() {
    //here we will do some sw/hw checks and return 1 if all ok, else 0
    //i2cdetect -l
    //i2c-1 i2c         bcm2835 I2C adapter               I2C adapter
    int ret_val = 0;

    if ((i2c_fd = open ("/dev/i2c-1", O_RDWR)) < 0)
        printf("Unable to open I2C device\n");

    if (ioctl (i2c_fd, I2C_SLAVE, I2C_DC_IC_ADR) < 0)
        printf("Unable to select I2C device:\n");

    if (i2c_fd == -1)
        printf("Error 001\n");
    else
        //cout << "I2C device opened OK with result " << fd << endl;
        ret_val = i2c_fd;

    driver_initialized = 1;

    memset(read_buff, 0x00, sizeof(read_buff));
    memset(write_buff, 0x00, sizeof(write_buff));

    return ret_val;
}

//---------------------------------------------------------------------------
/// @internal
///
/// Calculate the CRC8 of the byte value provided with the current
/// global 'crc8' value.
///
/// @param[in] data
/// data to compute crc8 on
///
/// @return
/// CRC8 result
/// @endinternal
///
unsigned char calc_crc8(unsigned char data)
{
    int i;
    // See Application Note 27
    crc8 = crc8 ^ data;
    for (i = 0; i < 8; ++i)
    {
        if (crc8 & 1)
            crc8 = (crc8 >> 1) ^ 0x8c;
        else
            crc8 = (crc8 >> 1);
    }
    return crc8;
}

//byte array to HEX string
char hexmap[16] = {'0', '1', '2', '3', '4', '5', '6', '7', '8', '9', 'A', 'B', 'C', 'D', 'E', 'F'};
char * hexStr2(unsigned char * data, int len) {
    //memset(buffTMP, 0x00, sizeof(buffTMP));
    for (int i = 0; i < len; ++i) {
        buffTMP[2 * i]     = hexmap[(data[i] & 0xF0) >> 4];
        buffTMP[2 * i + 1] = hexmap[data[i] & 0x0F];
    }
    buffTMP[2 * len] = 0x00;
    return buffTMP;
}
char * hexStr3(unsigned char * data, int len) {
    //memset(buffTMP, 0x00, sizeof(buffTMP));
    for (int i = 0; i < len; ++i) {
        buffTMP[3 * i]     = hexmap[(data[i] & 0xF0) >> 4];
        buffTMP[3 * i + 1] = hexmap[data[i] & 0x0F];
        buffTMP[3 * i + 2] = ' ';
    }
    buffTMP[3 * len] = 0x00;
    return buffTMP;
}

//GET unique ID from DS device
unsigned char * getDeepCoverID (int verbose) {

    if (driver_initialized == 0) initDriver();

    getPageData(PG_ROM_OPTIONS, verbose);

    //validate CRC-8 of ID
    crc8 = 0;
    if (verbose) printf("\n");
    for (int x = 0; x < 8; x++)  {
        calc_crc8(read_buff[26 + x]);
        //cout << "crc-8[" << x << "]=" << n2hexstr(crc8) << "\n";
    }

    // verify CRC and correct family code
    if (crc8 != 0)  printf("ERROR: CRC-8 for returned ID failed!\n");
    else if (((read_buff[26] & 0x7F) != 0x4C)) printf("ERROR: Not walid DS family code!\n");
    else if (verbose) printf("CRC-8 for returned ID is VALID!\n");

    return read_buff+26; //return adress of 26th position in our buffer
}

//Read RNG
unsigned char * getRNGdata(int len, int skip_header) {

    if (driver_initialized == 0) initDriver();

    write_buff[0] = CMD_READ_RNG;
    write_buff[1] = 0x01;
    write_buff[2] = len; //page we want to read (or len for CMD_READ_RNG)
    int out_count = write (i2c_fd, write_buff, 3);

    usleep(2000);

    //read reply now
    int redoutCount = read(i2c_fd, read_buff, 65); //

    if (skip_header) {
        return read_buff + 1;
    } else {
        return read_buff;
    }
}

//---------------------------------------------------------------------------
//-------- DS28C36 Memory functions
//---------------------------------------------------------------------------

//---------------------------------------------------------------------------
//Read page data from device
int getPageProtection(int page) {

    if (page >=0 && page < 32) {
	    if (driver_initialized == 0) initDriver();

	    /*
	    •	<Start, device address write>
	    •	TX: Read Memory Command
	    •	TX: Length (SMBus)
	    •	TX: Parameter
	    •	<Stop>
	    •	<Delay>
	    •	<Start, device address read>
	    •	RX: Length (SMBus)
	    •	RX: Result byte
	    •	RX: Data
	    •	<Stop>
	    */

	    write_buff[0] = CMD_READ_PAGE_PROT;
	    write_buff[1] = 0x01; //len
	    write_buff[2] = page; //page we want to read
	    int out_count = write (i2c_fd, write_buff, 3);

	    usleep(500);

	    //read reply now
	    int redoutCount = read(i2c_fd, read_buff, 2); //
	    return read_buff[1];

	} else {
		return 0;
	}
}

//---------------------------------------------------------------------------
//Read page data from device
unsigned char * getPageData(int page, int skip_header) {

    if (page >=0 && page < 32) {
	    if (driver_initialized == 0) initDriver();

	    /*
	    •	<Start, device address write>
	    •	TX: Read Memory Command
	    •	TX: Length (SMBus)
	    •	TX: Parameter
	    •	<Stop>
	    •	<Delay>
	    •	<Start, device address read>
	    •	RX: Length (SMBus)
	    •	RX: Result byte
	    •	RX: Data
	    •	<Stop>
	    */

	    write_buff[0] = CMD_READ_MEM;
	    write_buff[1] = 0x01; //len
	    write_buff[2] = page; //page we want to read
	    int out_count = write (i2c_fd, write_buff, 3);

	    usleep(1500);

	    //read reply now
	    int redoutCount = read(i2c_fd, read_buff, 32 + 2); //

	    if (skip_header) {
	        return read_buff + 2;
	    } else {
	        return read_buff;
	    }
	} else {
		memset(read_buff, 0x00, sizeof(read_buff));
		return read_buff;
	}
}

//---------------------------------------------------------------------------
//Write page data into DS device
int writePageData(int page, unsigned char * data) {
	if (page >=0 && page < 32) {
	    if (driver_initialized == 0) initDriver();

	    /*
	    •	<Start, device address write>
	    •	TX: Write Memory Command
	    •	TX: Length (SMBus) [always 33]
	    •	TX: Parameter
	    •	TX: Data
	    •	<Stop>
	    •	<Delay>
	    •	<Start, device address read>
	    •	RX: Length (SMBus) [always 1]
	    •	RX: Result byte
	    •	<Stop>
	    */

	    write_buff[0] = CMD_WRITE_MEM ;
	    write_buff[1] = 33;
	    write_buff[2] = page; //page we want to write

	    // fill the rest of write buffer
	    for (int i=0; i<=32; i++) write_buff [i+3] = data[i];

	    //write now
	    int out_count = write (i2c_fd, write_buff, 35);

	    //wait ...
	    usleep(12500); //critical value, EEPROM write is slow

        //read reply now
		int redoutCount = read(i2c_fd, read_buff, 2); //len + res_code
		return read_buff[1]; //Result byte

	}

}

//---------------------------------------------------------------------------
//Read buffer data from device
unsigned char * getBufferData(int len, int skip_header) {
    /*
    •	<Start, device address write>
    •	TX: Read Buffer Command
    •	<Stop>
    •	<Start, device address read>
    •	RX: Length (SMBus)
    •	RX: Data
    •	<Stop>
    */
    write_buff[0] = CMD_READ_BUF;
    int out_count = write (i2c_fd, write_buff, 1);

    //read reply now
    if (len > 80 || len==0) len=80;
    int redoutCount = read(i2c_fd, read_buff, len+1); //

    if (skip_header) {
        return read_buff + 1;
    } else {
        return read_buff;
    }

}

//---------------------------------------------------------------------------
//Write buffer data to device
int writeBufferData(int len, unsigned char * data) {
    /*
    •	<Start, device address write>
    •	TX: Write Buffer Command  (Parameter not used because duplication with SMBus length)
    •	TX: Length (SMBus) [length]
    •	TX: Buffer Data
    •	<Stop>
    */
    write_buff[0] = CMD_WRITE_BUF;
    write_buff[1] = len;

    // fill the rest of write buffer
    for (int i=0; i<=len; i++) write_buff [i+2] = data[i];

    return write (i2c_fd, write_buff, len+2);
}

//---------------------------------------------------------------------------
// ECDSA signature generation
//---------------------------------------------------------------------------
unsigned char * computeReadPageAuthentication(unsigned char * data32, int skip_header) {

    if (driver_initialized == 0) initDriver();

    // Preload challenge (data) in DS Buffer
    writeBufferData(32, data32);

    int pg = 0;
    int at = AT_ECDSA_KEYA;

    /*
    •	[Preload challenge (data) in Buffer]
    •	<Start, device address write>
    •	TX: Compute and Read Page Authentication command
    •	TX: Length (SMBus)
    •	TX: Parameter
    •	<Stop>
    •	<Delay>
    •	<Start, device address read>
    •	RX: Length (SMBus)
    •	RX: Result byte
    •	RX: HMAC (32 bytes) or Signature (64 bytes)
    •	<Stop>
    */

    //at = AT_ECDSA_KEYA; //  AT_ECDSA_KEYA, AT_HMAC_SECRETA

    write_buff[0] = CMD_COMP_READ_AUTH; //CMD_COMP_READ_AUTH
    write_buff[1] = 1;
    write_buff[2] = ((at & 0x07) << 5) | (pg & 0x1F); //pack params

    //write now
    int redoutCountW = write (i2c_fd, write_buff, 3);

    //wait ...
    usleep(40000); //critical value, ECDSA computation is slow

    //read reply now
    int redoutCount = read(i2c_fd, read_buff, 64+2); //len AT_ECDSA_KEYA=64, AT_HMAC_SECRETA=32

    if (skip_header) {
        return read_buff + 2;
    } else {
        return read_buff;
    }
}



