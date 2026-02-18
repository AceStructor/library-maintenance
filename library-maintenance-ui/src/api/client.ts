import axios from "axios";

export const albumApi = axios.create({
  baseURL: "http://zelda:5000",
  headers: {
    "Content-Type": "application/json",
  },
});

export const youtubeApi = axios.create({
  baseURL: "http://localhost:8080", 
  headers: {
    "Content-Type": "application/json",
  },
});
