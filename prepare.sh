# Make sure osqueryctl is stopped:
osqueryctl stop

echo "Creating extensions dir..."
mkdir -p /etc/osquery/extensions/

echo "Adding logging extension to osquery..."
cp log_ext/bin/log.ext /etc/osquery/extensions/
echo "/etc/osquery/extensions/log.ext" > /etc/osquery/extensions.load

echo "Adding table extension to osquery..."
cp table_ext/bin/foobar.ext /etc/osquery/extensions/
echo "/etc/osquery/extensions/foobar.ext" >> /etc/osquery/extensions.load

echo "Adding osquery config and flag files..."
cp osquery.flags /etc/osquery
cp osquery.conf /etc/osquery
