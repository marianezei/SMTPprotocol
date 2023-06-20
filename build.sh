#!/bin/bash

# Executa server-exc.sh
echo "Executando o primeiro script - servidor"
chmod +x server-exc.sh

# Executa client-exc.sh
echo "Executando o segundo - cliente"
chmod +x client-exc.sh

./server-exc.sh&
./client-exc.sh&


~                     
