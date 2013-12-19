package main

import (
  "encoding/binary"
  "fmt"
  "github.com/docopt/docopt.go"
  "io/ioutil"
  "os"
  "strconv"
)

type Header struct {
  validationWord    uint32
  version           uint8
  flags             uint8
  length            uint16
  lengthBytes       uint16
  spare             uint16
  simpleChecksum    uint16
  crcChecksumOffset uint16
  crcChecksum       uint32
}

var VERSION uint8 = 0
var VALIDATION_WORD uint32 = 825250625

/* Alignement is in KB*/
var VALID_ALIGNMENTS = []int{64, 128, 256, 512, 1024, 2048, 4096}
var ALIGNMENT int
var ALIGNMENT_BYTES int

var header Header

var CRC_TABLE = []int{ 0, 79764919, 159529838, 222504665, 319059676, 398814059,
445009330, 507990021, 638119352, 583659535, 797628118, 726387553, 890018660,
835552979, 1015980042, 944750013, 1276238704, 1221641927, 1167319070,
1095957929, 1595256236, 1540665371, 1452775106, 1381403509, 1780037320,
1859660671, 1671105958, 1733955601, 2031960084, 2111593891, 1889500026,
1952343757, -1742489888, -1662866601, -1851683442, -1788833735, -1960329156,
-1880695413, -2103051438, -2040207643, -1104454824, -1159051537, -1213636554,
-1284997759, -1389417084, -1444007885, -1532160278, -1603531939, -734892656,
-789352409, -575645954, -646886583, -952755380, -1007220997, -827056094,
-898286187, -231047128, -151282273, -71779514, -8804623, -515967244,
-436212925, -390279782, -327299027, 881225847, 809987520, 1023691545,
969234094, 662832811, 591600412, 771767749, 717299826, 311336399, 374308984,
453813921, 533576470, 25881363, 88864420, 134795389, 214552010, 2023205639,
2086057648, 1897238633, 1976864222, 1804852699, 1867694188, 1645340341,
1724971778, 1587496639, 1516133128, 1461550545, 1406951526, 1302016099,
1230646740, 1142491917, 1087903418, -1398421865, -1469785312, -1524105735,
-1578704818, -1079922613, -1151291908, -1239184603, -1293773166, -1968362705,
-1905510760, -2094067647, -2014441994, -1716953613, -1654112188, -1876203875,
-1796572374, -525066777, -462094256, -382327159, -302564546, -206542021,
-143559028, -97365931, -17609246, -960696225, -1031934488, -817968335,
-872425850, -709327229, -780559564, -600130067, -654598054, 1762451694,
1842216281, 1619975040, 1682949687, 2047383090, 2127137669, 1938468188,
2001449195, 1325665622, 1271206113, 1183200824, 1111960463, 1543535498,
1489069629, 1434599652, 1363369299, 622672798, 568075817, 748617968, 677256519,
907627842, 853037301, 1067152940, 995781531, 51762726, 131386257, 177728840,
240578815, 269590778, 349224269, 429104020, 491947555, -248556018, -168932423,
-122852000, -60002089, -500490030, -420856475, -341238852, -278395381,
-685261898, -739858943, -559578920, -630940305, -1004286614, -1058877219,
-845023740, -916395085, -1119974018, -1174433591, -1262701040, -1333941337,
-1371866206, -1426332139, -1481064244, -1552294533, -1690935098, -1611170447,
-1833673816, -1770699233, -2009983462, -1930228819, -2119160460, -2056179517,
1569362073, 1498123566, 1409854455, 1355396672, 1317987909, 1246755826,
1192025387, 1137557660, 2072149281, 2135122070, 1912620623, 1992383480,
1753615357, 1816598090, 1627664531, 1707420964, 295390185, 358241886,
404320391, 483945776, 43990325, 106832002, 186451547, 266083308, 932423249,
861060070, 1041341759, 986742920, 613929101, 542559546, 756411363, 701822548,
-978770311, -1050133554, -869589737, -924188512, -693284699, -764654318,
-550540341, -605129092, -475935807, -413084042, -366743377, -287118056,
-257573603, -194731862, -114850189, -35218492, -1984365303, -1921392450,
-2143631769, -2063868976, -1698919467, -1635936670, -1824608069, -1744851700,
-1347415887, -1418654458, -1506661409, -1561119128, -1129027987, -1200260134,
-1254728445, -1309196108 }

func intInSlice(a int, list []int) bool {
    for _, b := range list {
        if b == a {
            return true
        }
    }
    return false
}

func computeSimpleChecksum(checksum_data []byte) uint16 {
    var checksum uint16
    for _, data := range checksum_data {
      checksum += uint16(data & 0xFF)
    }
    return checksum
}

func computeCRCChecksum(checksum_data []byte) uint32{
    var checksum int
    var crc int = -1
    for _, data := range checksum_data {
        crc = crc << 8 ^ CRC_TABLE[((crc >> 24 ^ int(data)) & 0xFF)]
    }
    checksum = int(crc ^ 0xFFFFFFFF)
    return uint32(checksum & 0xFFFFFFFF)
}

func decodeHeader(source_ba []byte) {
  header.validationWord = binary.LittleEndian.Uint32(source_ba[64:68])
  invalid_word := ""
  if header.validationWord != VALIDATION_WORD {
    invalid_word = " (invalid)"
  }
  fmt.Printf(" Validation word: 0x%08x%s\n", header.validationWord, invalid_word)

  header.version = source_ba[68]
  invalid_version := ""
  if header.version != VERSION {
    invalid_version = " (invalid)"
  }else{
    fmt.Printf(" Version: %d%s\n", header.version, invalid_version)
  }

  header.flags = source_ba[69]
  fmt.Printf(" Flags: 0x%02x\n", header.flags)

  header.length = binary.LittleEndian.Uint16(source_ba[70:72])
  header.lengthBytes = header.length*4
  invalid_length := ""
  if header.length < 21 && header.length > 16384 {
    invalid_length = " (invalid)"
  }
  fmt.Printf(" Length: %d words (%d bytes)%s\n", header.length,
          header.lengthBytes, invalid_length)

  header.spare = binary.LittleEndian.Uint16(source_ba[72:74])
  fmt.Printf(" Spare: 0x%04x\n", header.spare)

  header.simpleChecksum = binary.LittleEndian.Uint16(source_ba[74:76])
  invalid_checksum := ""
  if header.simpleChecksum != computeSimpleChecksum(source_ba[64:74]) {
    invalid_checksum = " (invalid)"
  }
  fmt.Printf(" Simple Checksum: 0x%04x%s\n", header.simpleChecksum, invalid_checksum)

  if int(header.lengthBytes) <= len(source_ba) && header.lengthBytes > 0 {
    header.crcChecksumOffset = header.lengthBytes - 4
    checksum := source_ba[header.crcChecksumOffset:header.lengthBytes]
    header.crcChecksum = binary.LittleEndian.Uint32(checksum)
    invalid_checksum = ""
    if header.crcChecksum != computeCRCChecksum(source_ba[0:header.crcChecksumOffset]){
      invalid_checksum = " (invalid)"
    }
    fmt.Printf(" CRC Checksum: 0x%08x%s\n", header.crcChecksum, invalid_checksum)
  }
}

func createHeader(data []byte, size int) []byte {
  binary.LittleEndian.PutUint32(data[64:68], VALIDATION_WORD)
  data[68] = 0
  data[69] = 0
  /* file size is in bytes so we must / 4*/
  size = (size / 4) + 4 / 4
  binary.LittleEndian.PutUint16(data[70:72], uint16(size))
  data[72] = 0
  binary.LittleEndian.PutUint16(data[74:76], computeSimpleChecksum(data[64:74]))
  checksum_offset := (size * 4) - 4
  binary.LittleEndian.PutUint32(data[checksum_offset:checksum_offset + 4],computeCRCChecksum(data[0:checksum_offset]))
  return data
}

func main() {
  usage := `
usage:
  mkpimage [options] <source>...

options:
  -h, --help                  Display this help message and exit
  -v, --version=<version>     Header version to be created [default: 0]
  -a, --alignment=<alignment> Address alignment in kilobytes, valid value starts
                              from 64, 128, 256 etc, override if the NAND flash
                              has a larger block size [default: 64]
  -d, --decode                Flag to decode the header information from input file and display it
  -f, --force                 Flag to force decoding even if the input file is an unpadded image
  -o, --output=<destination>  Output file, relative and absolute path supported
`
    arguments, _ := docopt.Parse(usage, nil, true, "mkpimage", false)

    ALIGNMENT, _ = strconv.Atoi(arguments["--alignment"].(string))
    decode := arguments["--decode"].(bool)
    sources_path := arguments["<source>"].([]string)
    for _, source_path := range sources_path {
      if _, err := os.Stat(source_path); os.IsNotExist(err) {
        fmt.Printf("No such file or directory: %s", source_path)
        return
      }
    }

    /*Check that ALIGNMENT is a valid one.*/
    if ! intInSlice(ALIGNMENT, VALID_ALIGNMENTS) {
      fmt.Printf("%d is not a valid alignements, only %s supported.\n",
            ALIGNMENT, VALID_ALIGNMENTS)
    }

    ALIGNMENT_BYTES = ALIGNMENT * 1024

    if decode {
      source, err := os.Open(sources_path[0])
      if err != nil {
          return
      }
      defer source.Close()

      // get the file size
      stat, err := source.Stat()
      if err != nil {
          return
      }
      file_size := stat.Size()

      source_ba, _ := ioutil.ReadFile(sources_path[0])

      if int(file_size) == ALIGNMENT_BYTES {
        fmt.Println("Input file contains a single image:")
        decodeHeader(source_ba)
      }else{
        fmt.Println("Input file contains four images:")
        for idx := 0; idx < 4;  {
          fmt.Printf("Image: %d\n", idx)
          start_pos := idx * ALIGNMENT_BYTES
          end_pos := (idx + 1) * ALIGNMENT_BYTES
          decodeHeader(source_ba[start_pos:end_pos])
        }
      }
      return
    }

    /*Encode source into dest with a header*/
    if arguments["--output"] == nil {
      fmt.Println("You must give an output destination.")
    }
    output := arguments["--output"].(string)

    if len(sources_path) == 1 {
      source, err := os.Open(sources_path[0])
      if err != nil {
          return
      }
      defer source.Close()

      source_data, _ := ioutil.ReadFile(sources_path[0])
      data := make([]byte, ALIGNMENT_BYTES)
      copy(data, source_data)

      // get the file size
      stat, err := source.Stat()
      if err != nil {
          return
      }
      file_size := stat.Size()

      createHeader(data, int(file_size))
      ioutil.WriteFile(output, data, 0644)
    }else if len(sources_path) == 4 {
      /*merge 4 files into 1*/
      data := make([]byte, 4*ALIGNMENT_BYTES)
      var file_size int64
      for _, source_path := range sources_path {
        source, err := os.Open(source_path)
        if err != nil {
            return
        }
        defer source.Close()

        // get the file size
        stat, err := source.Stat()
        if err != nil {
            return
        }
        file_size += stat.Size()

        source_data, _ := ioutil.ReadFile(source_path)
        copy(data, source_data)
      }
      createHeader(data, int(file_size))
      ioutil.WriteFile(output, data, 0644)
    }
}

