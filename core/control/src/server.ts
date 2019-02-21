import express = require('express');
import { WebApi } from './application';

const port: number = Number(process.env.PORT) || 5001; // or from a configuration file
const api = new WebApi(express(), port);
api.run();
console.info(`listening on ${port}`);
