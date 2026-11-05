# Restreamer v1 database import

This programe reads a Restreamer v1 database (v1.json, created with 0.6.7+) and writes a new database including the metadata for the UI.
Also static routes will be added to the configuration file in order to main the paths to the player, snapshot, and m3u8.

This will make the migration from the old to the new Restreamer seamless.

If the v1 database is from a Restreamer before 0.6.7, then some assumptions are made that the input stream has a H264 video track.

All the environment variables for defining the USB device or Rapsicam are respected and are used to define the inputs.

## Processing

The importer understands and respects the same environment variables as the Restreamer. It uses `RS_CONFIGFILE` (and possibly
`RS_DB_DIR`) to find out where to look for the databases. If the v1 database is not found in this directory, there's nothing
to import and the importer aborts. If there's a v1 database, the importer checks if there's also a non-empty new database. If
there is an existing non-empty new database then the v1 database will not be imported. The existing non-empty new database is
interpreted that this is not a fresh Restreamer installation and should not be altered.

In case the new database is considered empty (i.e. no process and no metadata entries), then the v1 database will be read and
transformed to the corresponding new database format. In order keep compatibility with the v1 Restreamer, the importer also
alters the configuration file by adding static redirect routes to the new locations of the player, snapshot, and m3u8.

## Return Values

The importer exits with `0` if the import was successful or if the import did not happen (either because there wasn't a v1 database,
or there was already a non-empty new database). If there were any error during the importing process, the importer exits with `1`.
During the runtime of the importer, log messages are written to `stderr`.
