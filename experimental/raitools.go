package main

import (
	"fmt"
	ba "github.com/Workiva/go-datastructures/bitarray"
	"bytes"
	"strconv"
	"golang.org/x/crypto/blake2b"
)

/*
	Experimental Rai functions
	Needs tests and optimization
	Maybe change bitarray library
	Inspired by https://github.com/icarusglider/PyRai/blob/master/pyrai.py
 */

func XrbToAccount(address string) (string, error) {
	// Given a string containing an XRB address, confirm validity and provide resulting hex address
	if len(address) == 64 && (address[:4] == "xrb_") {
		// each index = binary value, account_lookup[0] == '1'
		account_map := "13456789abcdefghijkmnopqrstuwxyz"
		// extract after 'xrb_' and before the 8-char checksum
		acrop_key := address[4:len(address)-8]
		// extract checksum
		acrop_check := address[len(address)-8:]
		account_lookup := make(map[string]ba.BitArray, len(account_map))

		// populate lookup index with prebuilt bitarrays ready to append
		ba_size := 5
		for k, v := range account_map {
			// New bit array with length 5
			bval := ba.NewBitArray(uint64(ba_size))
			bi := fmt.Sprintf("%b", k)
			shift := ba_size - len(bi)
			for i, j := range bi {
				// max len 5: 0 to 4, but bi len could be less, so shift = 5 - bi_len
				if j == 49 {
					bval.SetBit(uint64(shift + i))
				}
			}
			account_lookup[string(v)] = bval
		}
		number_size := 256
		number_l := ba.NewBitArray(uint64(number_size))
		first := true
		i := 0
		for _, v := range acrop_key {
			ltr := string(v)
			tmp_ba := account_lookup[ltr]
			if first {
				// Only get fifth bit for the first BitArray
				bval, _ := tmp_ba.GetBit(uint64(4))
				if bval {
					number_l.SetBit(uint64(i))
				}
				first = false
				i += 1
			} else {
				for j := 0; j < 5; j++ {
					bval, _ := tmp_ba.GetBit(uint64(j))
					if bval {
						number_l.SetBit(uint64(i))
					}
					i += 1
				}
			}
		}
		check_size := len(acrop_check) * ba_size
		check_l := ba.NewBitArray(uint64(check_size))
		i = 0
		for _, v := range acrop_check {
			ltr := string(v)
			tmp_ba := account_lookup[ltr]
			for j := 0; j < 5; j++ {
				bval, _ := tmp_ba.GetBit(uint64(j))
				if bval {
					check_l.SetBit(uint64(i))
				}
				i += 1
			}
		}
		// Prepare result
		result := ToHexString(number_l, number_size)
		check_swap := ByteSwap(check_l, check_size)
		numberArr := ToByteArray(number_l, number_size)
		digest, err := blake2b.New(ba_size, []byte{})
		if err != nil {
			fmt.Print(err)
			return "", err
		}
		digest.Write(numberArr)
		d := digest.Sum(nil)
		hexdigest := fmt.Sprintf("%X", d)
		if hexdigest == ToHexString(check_swap, check_size) {
			return result, nil
		} else {
			return "", nil
		}
	}
	return "", nil
}

func AccountToXrb() {

}

/*
def account_xrb(account):
	# Given a string containing a hex address, encode to public address format with checksum
	# each index = binary value, account_lookup['00001'] == '3'
	account_map = "13456789abcdefghijkmnopqrstuwxyz"
	account_lookup = {}
	# populate lookup index for binary string to base-32 string character
	for i in range(0,32):
		account_lookup[BitArray(uint=i,length=5).bin] = account_map[i]
	# hex string > binary
	account = BitArray(hex=account)

	# get checksum
	h = blake2b(digest_size=5)
	h.update(account.bytes)
	checksum = BitArray(hex=h.hexdigest())

	# encode checksum
	# swap bytes for compatibility with original implementation
	checksum.byteswap()
	encode_check = ''
	for x in range(0,int(len(checksum.bin)/5)):
		# each 5-bit sequence = a base-32 character from account_map
		encode_check += account_lookup[checksum.bin[x*5:x*5+5]]

	# encode account
	encode_account = ''
	# pad our binary value so it is 260 bits long before conversion (first value can only be 00000 '1' or 00001 '3')
	while len(account.bin) < 260:
		account = '0b0' + account
	for x in range(0,int(len(account.bin)/5)):
		# each 5-bit sequence = a base-32 character from account_map
		encode_account += account_lookup[account.bin[x*5:x*5+5]]
	# build final address string
	return 'xrb_'+encode_account+encode_check
 */

func ToBitString(bits ba.BitArray, length int) string {
	var buffer bytes.Buffer
	for i := 0; i < length; i++ {
		bval, _ := bits.GetBit(uint64(i))
		if bval {
			buffer.WriteString("1")
		} else {
			buffer.WriteString("0")
		}
	}
	return buffer.String()
}

func ToHexString(bits ba.BitArray, length int) string {
	hexarray := []string{"0", "1", "2", "3", "4", "5", "6", "7", "8", "9", "A", "B", "C", "D", "E", "F"}
	//barray := make([]byte, length/4, length/4)
	var tmpBuffer bytes.Buffer
	var hexBuffer bytes.Buffer
	j, k := 0, 0
	for i := 0; i < length; i++ {
		if j == 4 {
			// get 4 bit str
			intr, err := strconv.ParseInt(tmpBuffer.String(), 2, 64)
			tmpBuffer.Reset()
			if err != nil {
				fmt.Println(err)
			}
			hexBuffer.WriteString(hexarray[intr])
			//barray[k] = byte(intr)
			j = 0
			k++
		}
		bval, _ := bits.GetBit(uint64(i))
		if bval {
			tmpBuffer.WriteString("1")
		} else {
			tmpBuffer.WriteString("0")
		}
		j++
	}
	// get last 4 bit str
	intr, err := strconv.ParseInt(tmpBuffer.String(), 2, 64)
	if err != nil {
		fmt.Println(err)
	}
	hexBuffer.WriteString(hexarray[intr])
	//barray[k] = byte(intr)
	return hexBuffer.String()
}

func ToByteArray(bits ba.BitArray, length int) []byte {
	barray := make([]byte, length/8, length/8)
	var tmpBuffer bytes.Buffer
	j, k := 0, 0
	for i := 0; i < length; i++ {
		if j == 8 {
			// get 8 bit str
			intr, err := strconv.ParseInt(tmpBuffer.String(), 2, 64)
			tmpBuffer.Reset()
			if err != nil {
				fmt.Println(err)
			}
			barray[k] = byte(intr)
			j = 0
			k++
		}
		bval, _ := bits.GetBit(uint64(i))
		if bval {
			tmpBuffer.WriteString("1")
		} else {
			tmpBuffer.WriteString("0")
		}
		j++
	}
	// get last 8 bit str
	intr, err := strconv.ParseInt(tmpBuffer.String(), 2, 64)
	if err != nil {
		fmt.Println(err)
	}
	barray[k] = byte(intr)
	return barray
}

func ByteSwap(src ba.BitArray, length int) ba.BitArray {
	out := ba.NewBitArray(uint64(length))
	wordLen := 8
	j := 0
	cnt := 1
	for i := 0; i < length; i++ {
		if j == wordLen {
			j = 0
			cnt++
		}
		bval, _ := src.GetBit(uint64(i))
		if bval {
			// Get the right insert position, j is the delta offset
			pos := (length - cnt*wordLen) + j
			out.SetBit(uint64(pos))
		}
		j++
	}
	return out
}

func main() {
	address := "xrb_1ipx847tk8o46pwxt5qjdbncjqcbwcc1rrmqnkztrfjy5k7z4imsrata9est"
	fmt.Printf("Address %s\n", address)

	// Given a string containing an XRB address, confirm validity and
	// provide resulting hex address
	res, err := XrbToAccount(address)
	if err != nil {
		fmt.Printf("%s\n", err)
	}
	fmt.Printf("Resulting hex address %s\n", res)

	fmt.Printf("Done\n")
}
