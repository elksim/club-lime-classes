console.log(`settings.js run!`);

/** @type {[Set<string>, Set<string>]} */
let [selectedWorkouts, selectedLocations] = [new Set(), new Set()];

function saveSetToLocalStorage(key, set) {
    localStorage.setItem(key, JSON.stringify([...set]));
}

function getSetFromLocalStorage(key) {
    return new Set(JSON.parse(localStorage.getItem(key) || "[]"));
}

document.addEventListener("DOMContentLoaded", () => {
    console.log(`DOMContentLoaded`);

    selectedLocations = getSetFromLocalStorage("selectedLocations");
    selectedWorkouts = getSetFromLocalStorage("selectedWorkouts");

    populateLocations();
    populateWorkouts();
})


//region: locations
function populateLocations() {
    const tmp = /** @type {HTMLSelectElement} */ document.querySelector("#selectedState");
    // @ts-ignore
    const selectedState = tmp.value;
    if (selectedState === "") {
        return
    }

    //populate locations Element
    /** @type {HTMLElement} */
    // @ts-ignore
    let locationsEl = document.querySelector("#locations");
    locationsEl.innerHTML = "";

    // @ts-ignore
    for (const newLocation of stateToLocations[selectedState]) {
        /** @type HTMLElement */
        let newEl = document.createElement("div");
        newEl.classList.add("location");
        newEl.onclick = handleLocationClick;
        newEl.innerHTML = newLocation;
        if (selectedLocations.has(newLocation)) {
            newEl.classList.add("selected");
        }
        locationsEl.appendChild(newEl);
    }

    toggleLocationsOn = true;
    let toggleAllLocationsEl = document.querySelector("#toggleAllLocations");
    // @ts-ignore 
    toggleAllLocationsEl.innerHTML = "select entire state";
}

function handleSelectedStateChange() {
    populateLocations()
}
/** @param {MouseEvent} event */
function handleLocationClick(event) {
    let target = /** @type {HTMLElement} */ (event.target);
    if (target.classList.contains("selected")) {
        target.classList.remove("selected");
        selectedLocations.delete(target.innerHTML)
    } else {
        target.classList.add("selected");
        selectedLocations.add(target.innerHTML)
    }
    console.log(`selectedLocations:`, selectedLocations);
}

let toggleLocationsOn = true;
function handleToggleAllLocationsClick(event) {
    console.log(`xd ${event}`);
    console.log(`toggleOn: ${toggleLocationsOn}`);
    /** @type {HTMLElement} */
    //@ts-ignore
    const locationsEl = document.querySelector("#locations")

    if (toggleLocationsOn) {
        for (const child of locationsEl.children) {
            child.classList.add("selected")
            selectedLocations.add(child.innerHTML)
        }
    } else {
        for (const child of locationsEl.children) {
            child.classList.remove("selected")
            selectedLocations.delete(child.innerHTML)
        }
    }
    event.target.innerHTML = toggleLocationsOn ? "select entire state" : "deselect entire state";
    toggleLocationsOn = !toggleLocationsOn;
}

function handleSaveLocations(event) {
    const innerHTMLBefore = event.target.innerHTML;
    event.target.innerHTML = "saved locations!";
    setTimeout(() => {
        event.target.innerHTML = innerHTMLBefore;
    }, 1000);
    saveSetToLocalStorage("selectedLocations", selectedLocations);
}


//region: workouts

let toggleWorkoutsOn = true;
function handleToggleAllWorkoutsClick(event) {
    let workoutsElement = document.querySelector("#workouts");
    if (confirm("are you sure you want to toggle all workouts")) {
        for (const child of workoutsElement.children) {
            if (toggleWorkoutsOn) {
                child.classList.add("selected");
                selectedWorkouts.add(child.innerHTML);
            } else {
                child.classList.remove("selected");
                selectedWorkouts.delete(child.innerHTML);
            }
        }
    }
    event.target.innerHTML = toggleWorkoutsOn ? "select all workouts" : "deselect all workouts";
    toggleWorkoutsOn = !toggleWorkoutsOn;
}

function populateWorkouts() {
    console.log(`xd`);
    console.log(`workouts: `, workouts);
    console.log(`locations: `, locations);
    let workoutsEl = document.querySelector("#workouts");
    for (const workout of workouts) {
        let newEl = document.createElement("div");
        newEl.classList.add("workout");
        newEl.onclick = handleWorkoutClick;
        newEl.innerHTML = workout;
        if (selectedWorkouts.has(workout)) {
            newEl.classList.add("selected");
        }
        workoutsEl.appendChild(newEl);
    }
}

/** @param {MouseEvent} event */
function handleWorkoutClick(event) {
    let target = /** @type {HTMLElement} */ (event.target);
    if (target.classList.contains("selected")) {
        target.classList.remove("selected");
        selectedWorkouts.delete(target.innerHTML)
    } else {
        target.classList.add("selected");
        selectedWorkouts.add(target.innerHTML)
    }
    console.log(`selectedWorkouts:`, selectedWorkouts);
}



function handleSaveWorkouts(event) {
    const innerHTMLBefore = event.target.innerHTML;
    event.target.innerHTML = "saved workouts!";
    setTimeout(() => {
        event.target.innerHTML = innerHTMLBefore;
    }, 1000);
    saveSetToLocalStorage("selectedWorkouts", selectedWorkouts);
}