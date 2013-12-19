mkpimage
========

mkpimage is a tool used to sign U-Boot SPL binary for SoCFPGA platform.

This tool was created after Altera's tool. The only reason it exists is because
Altera's tool is done in Java and has hardcoded path to its JRE...
This makes it pretty useless as it's really difficult to integrate into an
automatic build system.

Which forces you to download anh install Altera's IDE which you may not want if
you don't care about Quartus (tool to work with the FPGA).
