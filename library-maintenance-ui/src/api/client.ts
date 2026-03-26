import axios from "axios";

export const api = axios.create({
  baseURL: "http://zelda:5000",
  headers: {
    "Content-Type": "application/json",
  },
});

export const localApi = axios.create({
  baseURL: "http://zelda:5001", 
  headers: {
    "Content-Type": "application/json",
  },
});
