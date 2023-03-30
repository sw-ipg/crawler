select i.crc32_checksum
    from index i
    where i.path = $1