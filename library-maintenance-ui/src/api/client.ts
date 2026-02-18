import axios from "axios";

export const albumApi = axios.create({
  baseURL: "http://zelda:5000",
  headers: {
    "Content-Type": "application/json",
  },
});

export const youtubeApi = axios.create({
  baseURL: "http://zelda:6000", 
  headers: {
    "Content-Type": "application/json",
  },
});
