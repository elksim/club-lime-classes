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
    toggleAllLocationsEl.innerHTML = "select all";
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
    event.target.innerHTML = toggleLocationsOn ? "deselect all" : "select all";
    toggleLocationsOn = !toggleLocationsOn;
}

function handleSaveLocations(event) {
    console.log(`saving!`);
    console.log(`selectedLocations: `, selectedLocations);
    saveSetToLocalStorage("selectedLocations", selectedLocations);
}


//region: workouts

