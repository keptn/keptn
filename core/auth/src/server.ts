import express = require('express');
import { WebApi } from './application';

let port: number = parseInt(process.env.PORT) || 5001; //or from a configuration file
let api = new WebApi(express(), port);
api.run();
console.info(`listening on ${port}`);