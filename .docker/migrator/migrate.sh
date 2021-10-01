apt-get -y update
apt-get -y install curl
curl -X POST 'http://search:9308/sql' -d "mode=raw&query=CREATE TABLE IF NOT EXISTS products (articul text, shop text, title text) min_infix_len='10'"
exit