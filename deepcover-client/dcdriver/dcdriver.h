/******************************************************************************
/// @file       dcdriver.h
/// @date       22.11.2018
/// @copyright  Copyright (C) beyond protocol inc. - All Rights Reserved
///             You may use, distribute and modify this code under the
///             terms and conditions defined in the file 'LICENSE.txt',
///             which is part of this source code package.
/// @brief		Raspberry Pi DS28C36 I2C interface driver
/// @details
******************************************************************************/

#ifndef DCDRIVERH_C_H_
#define DCDRIVERH_C_H_

//low-level and test functions
int initDriver(void);
unsigned char * getRNGdata(int len, int skip_header);
unsigned char * getPageData(int page, int skip_header);
unsigned char * getBufferData(int len, int skip_header);
int writePageData(int page, unsigned char * data);
int writeBufferData(int len, unsigned char * data);
int getPageProtection(int page);

//high-level functions
unsigned char * getDeepCoverID (int verbose); //return 8 bytes
unsigned char * computeReadPageAuthentication(unsigned char * data32, int skip_header); //return 66 bytes (len+rcode+sigs32+sigr32) with header

//helper
char * hexStr2(unsigned char * data, int len);
char * hexStr3(unsigned char * data, int len);
char * printPageName (int p);

#endif  // DCDRIVERH_C_H_

