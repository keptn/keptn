#!/bin/bash
echo "Usage: ./add-to-cart.sh http://XX.XX.XX.XX/carts/1/items"
echo "Press [CTRL+C] to stop."

if [ -z $1 ]
then
    echo "Please provide the url as parameter"
    echo ""
    echo "Usage: ./add-to-cart.sh http://XX.XX.XX.XX/carts/1/items"
    exit 1
fi

url=$1

i=0
while true
do
  echo ""
  echo "adding item to cart..."
  curl -X POST -H "Content-Type: application/json" -d "{\"id\":\"3395a43e-2d88-40de-b95f-e00e1502085b\", \"itemId\":\"03fef6ac-1896-4ce8-bd69-b798f85c6e0b\"}" $url
  i=$((i+1))
  if [ $i -ge 100 ] 
  then
    curl -X DELETE http://35.232.216.57/carts/1
    i=0
  fi
done
