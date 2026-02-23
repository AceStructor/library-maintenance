import axios from "axios";

export const albumApi = axios.create({
  baseURL: "http://zelda:5000",
  headers: {
    "Content-Type": "application/json",
  },
});

export const songApi = axios.create({
  baseURL: "http://zelda:5001", 
  headers: {
    "Content-Type": "application/json",
  },
});
