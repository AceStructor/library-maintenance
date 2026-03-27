<template>
  <div>
    <h3>Create Playlist</h3>

    <!-- Input Section -->
    <div>
      <input
        v-model="name"
        placeholder="Namen eingeben (Optional)"
        type="text"
      />

      <div>
        <select
        v-model="selectedMonth"
        class="bg-gray-700 border border-gray-600 px-2 py-1 rounded"
        >
        <option disabled value="January">Monat auswählen</option>
        <option
            v-for="m in months"
            :key="m"
            :value="m"
        >
            {{ m }}
        </option>
        </select>

        <select
        v-model="selectedYear"
        class="bg-gray-700 border border-gray-600 px-2 py-1 rounded"
        >
        <option disabled value="26">Jahr auswählen</option>
        <option
            v-for="y in years"
            :key="y"
            :value="y"
        >
            {{ y }}
        </option>
        </select>
      </div>

      <label>
      <input type="checkbox" v-model="autoMode" />
      Playlist automatisch erstellen
      </label>

      <div v-if="!autoMode" class="flex gap-2 items-center">
        <select
        v-model="selectedInterval"
        class="bg-gray-700 border border-gray-600 px-2 py-1 rounded"
        >
        <option disabled value="1 Month">Intervall auswählen</option>
        <option
            v-for="i in intervals"
            :key="i"
            :value="i"
        >
            {{ i }}
        </option>
        </select>

        <select
        v-model="selectedWildness"
        class="bg-gray-700 border border-gray-600 px-2 py-1 rounded"
        >
        <option disabled value=0>Wildheit auswählen</option>
        <option
            v-for="w in [0,1,2,3]"
            :key="w"
            :value="w"
        >
            {{ w }}
        </option>
        </select>

        <input
          v-model="playlistLength"
          placeholder="Länge eingeben"
          type="number"
        />
      </div>



      <button @click="create" :disabled="creating">
        {{ creating ? "Erstelle..." : "Erstellen" }}
      </button>
    </div>

    <p v-if="error" style="color:red">{{ error }}</p>
  </div>
</template>

<script setup lang="ts">
import { ref } from "vue";
import { api } from "../api/client.js";

const months = ["January", "February", "March", "April", "May", "June", "July", "August", "September", "October", "November", "December"];
const years = ["26"];
const intervals = ["1 month", "2 months", "3 months", "6 months", "1 year"];
const selectedMonth = ref("January");
const selectedYear = ref("26");
const name = ref("");
const autoMode = ref(true);
const creating = ref(false);
const selectedWildness = ref(0);
const selectedInterval = ref("1 month");
const playlistLength = ref(50);

const error = ref("");

const create = async () => {
  creating.value = true;

  try {
    await api.post("/playlist/add", {
      "name": name.value,
      "month": selectedMonth.value,
      "year": selectedYear.value,
      "auto": autoMode.value,
      "wildness": selectedWildness.value,
      "interval": selectedInterval.value,
      "length": playlistLength.value
    });
  } catch (err: any) {
    error.value = err.message;
  } finally {
    creating.value = false;
  }
};
</script>

<style scoped>
table {
  width: 100%;
  border-collapse: collapse;
  margin-top: 15px;
}
td, th {
  border: 1px solid #ddd;
  padding: 8px;
}
tr:hover {
  background-color: #434547;
  cursor: pointer;
}
.selected {
  background-color: #643a3a;
}
.detail-section {
  margin-top: 20px;
}

.details input {
  margin-bottom: 0.5rem;
}
.details .checkbox-row {
  margin-top: 0.5rem;
}
</style>
