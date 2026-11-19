"use strict";

// ---------------------------------------------------------------------------------------------------------------------
// Machine definitions
// ---------------------------------------------------------------------------------------------------------------------

const ALPHABET = "ABCDEFGHIJKLMNOPQRSTUVWXYZ";
// The German QWERTZ layout used on Enigma keyboards, lampboards and plugboards.
const BOARD_ROWS = ["QWERTZUIO", "ASDFGHJK", "PYXCVBNML"];
const ALL_ROTORS = ["I", "II", "III", "IV", "V", "VI", "VII", "VIII"];
const SLOTS = ["Left", "Middle", "Right"];

const MODELS = {
  I: {
    label: "Enigma I",
    note: "Used by the Army and Air Force. Three rotors chosen from I–V.",
    rotors: ALL_ROTORS.slice(0, 5),
    reflectors: ["A", "B", "C"],
    fourth: false,
  },
  M3: {
    label: "M3",
    note: "Used by the Navy. Three rotors chosen from I–VIII.",
    rotors: ALL_ROTORS,
    reflectors: ["B", "C"],
    fourth: false,
  },
  M4: {
    label: "M4",
    note: "Used by U-boats. A fixed fourth rotor (β or γ) sits beside a thin reflector.",
    rotors: ALL_ROTORS,
    reflectors: ["B-thin", "C-thin"],
    fourth: true,
  },
};

const REFLECTOR_LABELS = { A: "UKW-A", B: "UKW-B", C: "UKW-C", "B-thin": "UKW-B thin", "C-thin": "UKW-C thin" };
const FOURTH_ROTORS = { Beta: "β", Gamma: "γ" };

const DEFAULT_STATE = {
  model: "I",
  rotors: ["I", "II", "III"],
  rings: ["A", "A", "A"],
  positions: ["A", "A", "A"],
  fourthRotor: "Beta",
  fourthPosition: "A",
  reflector: "B",
  pairs: [],
  text: "",
};

const STORAGE_KEY = "ultra:settings:v1";

// ---------------------------------------------------------------------------------------------------------------------
// State
// ---------------------------------------------------------------------------------------------------------------------

let state = loadState();
let engineReady = false;
let lastResult = null;
let pendingSocket = null;

function loadState() {
  try {
    const saved = JSON.parse(localStorage.getItem(STORAGE_KEY) || "null");
    if (saved && typeof saved === "object") {
      return normalise({ ...structuredClone(DEFAULT_STATE), ...saved });
    }
  } catch {
    // Storage can be unavailable (private windows, blocked cookies); fall back to defaults.
  }
  return structuredClone(DEFAULT_STATE);
}

function saveState() {
  try {
    localStorage.setItem(STORAGE_KEY, JSON.stringify(state));
  } catch {
    // Not being able to remember settings is fine.
  }
}

// Make sure the state describes a machine that really existed, e.g. after switching from an M3 to an Enigma I.
function normalise(s) {
  if (!MODELS[s.model]) s.model = DEFAULT_STATE.model;
  const model = MODELS[s.model];
  const isLetter = (l) => typeof l === "string" && l.length === 1 && ALPHABET.includes(l);

  const used = new Set();
  s.rotors = [0, 1, 2].map((i) => {
    const r = Array.isArray(s.rotors) ? s.rotors[i] : null;
    if (model.rotors.includes(r) && !used.has(r)) {
      used.add(r);
      return r;
    }
    return null;
  });
  s.rotors = s.rotors.map((r) => {
    if (r) return r;
    const pick = model.rotors.find((candidate) => !used.has(candidate));
    used.add(pick);
    return pick;
  });

  s.rings = [0, 1, 2].map((i) => (Array.isArray(s.rings) && isLetter(s.rings[i]) ? s.rings[i] : "A"));
  s.positions = [0, 1, 2].map((i) => (Array.isArray(s.positions) && isLetter(s.positions[i]) ? s.positions[i] : "A"));
  if (!FOURTH_ROTORS[s.fourthRotor]) s.fourthRotor = "Beta";
  if (!isLetter(s.fourthPosition)) s.fourthPosition = "A";
  if (!model.reflectors.includes(s.reflector)) s.reflector = model.reflectors.includes("B") ? "B" : "B-thin";

  const seen = new Set();
  s.pairs = (Array.isArray(s.pairs) ? s.pairs : []).filter((p) => {
    if (!Array.isArray(p) || p.length !== 2 || !isLetter(p[0]) || !isLetter(p[1]) || p[0] === p[1]) return false;
    if (seen.has(p[0]) || seen.has(p[1])) return false;
    seen.add(p[0]);
    seen.add(p[1]);
    return true;
  }).slice(0, 13);

  if (typeof s.text !== "string") s.text = "";
  return s;
}

function engineConfig() {
  const model = MODELS[state.model];
  return {
    rotors: state.rotors,
    positions: state.positions.join(""),
    rings: state.rings.join(""),
    fourthRotor: model.fourth ? state.fourthRotor : "",
    fourthPosition: state.fourthPosition,
    reflector: state.reflector,
    plugboard: state.pairs.map((p) => p.join("")).join(" "),
  };
}

const shift = (letter, by) => ALPHABET[(ALPHABET.indexOf(letter) + by + 26) % 26];

// ---------------------------------------------------------------------------------------------------------------------
// DOM helpers
// ---------------------------------------------------------------------------------------------------------------------

const $ = (id) => document.getElementById(id);

function el(tag, attrs = {}, ...children) {
  const node = document.createElement(tag);
  for (const [key, value] of Object.entries(attrs)) {
    if (key === "class") node.className = value;
    else if (key === "text") node.textContent = value;
    else if (key.startsWith("on")) node.addEventListener(key.slice(2), value);
    else node.setAttribute(key, value);
  }
  node.append(...children);
  return node;
}

function options(select, values, selected, label = (v) => v) {
  select.replaceChildren(...values.map((v) => el("option", { value: v, text: label(v) })));
  select.value = selected;
}

// ---------------------------------------------------------------------------------------------------------------------
// Settings: model, reflector and rotors
// ---------------------------------------------------------------------------------------------------------------------

function renderModelPicker() {
  const picker = $("model-picker");
  picker.replaceChildren(
    ...Object.entries(MODELS).map(([id, model]) =>
      el("button", {
        type: "button",
        role: "radio",
        "aria-checked": String(state.model === id),
        text: model.label,
        title: model.note,
        onclick: () => {
          state.model = id;
          normalise(state);
          renderSettings();
          update();
        },
      }),
    ),
  );
}

function renderReflector() {
  options($("reflector"), MODELS[state.model].reflectors, state.reflector, (r) => REFLECTOR_LABELS[r]);
}

function stepButtons(label, onStep, windowLetter) {
  return el(
    "div",
    { class: "rotor-window" },
    el("button", { type: "button", class: "step-btn", "aria-label": `${label}: next start letter`, text: "▲", onclick: () => onStep(1) }),
    windowLetter,
    el("button", { type: "button", class: "step-btn", "aria-label": `${label}: previous start letter`, text: "▼", onclick: () => onStep(-1) }),
  );
}

function renderRotorSlots() {
  const container = $("rotor-slots");
  const slots = [];
  const model = MODELS[state.model];

  if (model.fourth) {
    const select = el("select", { "aria-label": "Fourth rotor" });
    options(select, Object.keys(FOURTH_ROTORS), state.fourthRotor, (r) => `${FOURTH_ROTORS[r]} ${r}`);
    select.addEventListener("change", () => {
      state.fourthRotor = select.value;
      update();
    });

    const windowLetter = el("div", { class: "window-letter", id: "window-fourth", text: state.fourthPosition });
    slots.push(
      el(
        "div",
        { class: "slot" },
        el("span", { class: "label", text: "Fourth" }),
        select,
        stepButtons("Fourth rotor", (by) => {
          state.fourthPosition = shift(state.fourthPosition, by);
          renderStartNotes();
          update();
        }, windowLetter),
        el("span", { class: "start-note", id: "start-fourth" }),
        el("span", { class: "no-ring", text: "fixed" }),
      ),
    );
  }

  SLOTS.forEach((slotName, i) => {
    const select = el("select", { "aria-label": `${slotName} rotor` });
    options(select, model.rotors, state.rotors[i]);
    select.addEventListener("change", () => {
      // Each rotor exists only once, so picking one that is already fitted swaps the two over.
      const other = state.rotors.indexOf(select.value);
      if (other !== -1) state.rotors[other] = state.rotors[i];
      state.rotors[i] = select.value;
      renderRotorSlots();
      update();
    });

    const ring = el("select", { "aria-label": `${slotName} rotor ring setting` });
    options(ring, [...ALPHABET], state.rings[i], (l) => `${l} ${String(ALPHABET.indexOf(l) + 1).padStart(2, "0")}`);
    ring.addEventListener("change", () => {
      state.rings[i] = ring.value;
      update();
    });

    const windowLetter = el("div", { class: "window-letter", id: `window-${i}`, text: state.positions[i] });
    slots.push(
      el(
        "div",
        { class: "slot" },
        el("span", { class: "label", text: slotName }),
        select,
        stepButtons(`${slotName} rotor`, (by) => {
          state.positions[i] = shift(state.positions[i], by);
          renderStartNotes();
          update();
        }, windowLetter),
        el("span", { class: "start-note", id: `start-${i}` }),
        el("label", { class: "ring-select" }, "Ring", ring),
      ),
    );
  });

  container.replaceChildren(...slots);
  renderStartNotes();
}

function renderStartNotes() {
  state.positions.forEach((p, i) => ($(`start-${i}`).textContent = `start ${p}`));
  const fourth = $("start-fourth");
  if (fourth) fourth.textContent = `start ${state.fourthPosition}`;
}

function renderSettings() {
  renderModelPicker();
  renderReflector();
  renderRotorSlots();
}

$("reflector").addEventListener("change", (e) => {
  state.reflector = e.target.value;
  update();
});

// ---------------------------------------------------------------------------------------------------------------------
// Lampboard and keyboard
// ---------------------------------------------------------------------------------------------------------------------

function buildBoards() {
  $("lampboard").replaceChildren(
    ...BOARD_ROWS.map((row) =>
      el("div", { class: "board-row" }, ...[...row].map((l) => el("div", { class: "lamp", id: `lamp-${l}`, text: l }))),
    ),
  );
  $("keyboard").replaceChildren(
    ...BOARD_ROWS.map((row) =>
      el(
        "div",
        { class: "board-row" },
        ...[...row].map((l) =>
          el("button", { type: "button", class: "key", id: `key-${l}`, text: l, "aria-label": `Key ${l}`, onclick: () => pressKey(l) }),
        ),
      ),
    ),
  );
}

function flashKey(letter) {
  const key = $(`key-${letter}`);
  if (!key) return;
  key.classList.add("pressed");
  setTimeout(() => key.classList.remove("pressed"), 120);
}

function pressKey(letter) {
  state.text += letter;
  $("input-text").value = state.text;
  flashKey(letter);
  update();
}

function removeLastLetter() {
  state.text = state.text.replace(/[^A-Za-z]*[A-Za-z][^A-Za-z]*$/, "");
  $("input-text").value = state.text;
  update();
}

document.addEventListener("keydown", (e) => {
  if (e.ctrlKey || e.metaKey || e.altKey) return;
  if (e.target.closest("input, textarea, select")) return;

  if (/^[a-zA-Z]$/.test(e.key)) {
    e.preventDefault();
    pressKey(e.key.toUpperCase());
  } else if (e.key === "Backspace") {
    e.preventDefault();
    removeLastLetter();
  }
});

// ---------------------------------------------------------------------------------------------------------------------
// Message
// ---------------------------------------------------------------------------------------------------------------------

const groupsOfFive = (letters) => letters.match(/.{1,5}/g)?.join(" ") ?? "";

$("input-text").addEventListener("input", (e) => {
  state.text = e.target.value;
  update();
});

$("clear").addEventListener("click", () => {
  state.text = "";
  $("input-text").value = "";
  update();
});

$("swap").addEventListener("click", () => {
  if (!lastResult || !lastResult.output) return;
  state.text = groupsOfFive(lastResult.output);
  $("input-text").value = state.text;
  update();
});

$("copy-output").addEventListener("click", async (e) => {
  const button = e.currentTarget;
  try {
    await navigator.clipboard.writeText(groupsOfFive(lastResult?.output ?? ""));
    button.textContent = "Copied";
  } catch {
    button.textContent = "Copy failed";
  }
  setTimeout(() => (button.textContent = "Copy"), 1200);
});

// ---------------------------------------------------------------------------------------------------------------------
// Plugboard
// ---------------------------------------------------------------------------------------------------------------------


function buildPlugboard() {
  $("plugboard").replaceChildren(
    ...BOARD_ROWS.map((row) =>
      el(
        "div",
        { class: "board-row" },
        ...[...row].map((l) =>
          el(
            "button",
            { type: "button", class: "socket", id: `socket-${l}`, onclick: () => socketClick(l) },
            el("span", { class: "letter", text: l }),
            el("span", { class: "holes" }),
          ),
        ),
      ),
    ),
  );
}

function socketClick(letter) {
  const pairIndex = state.pairs.findIndex((p) => p.includes(letter));
  if (pairIndex !== -1) {
    state.pairs.splice(pairIndex, 1);
    pendingSocket = null;
  } else if (pendingSocket === letter) {
    pendingSocket = null;
  } else if (pendingSocket) {
    if (state.pairs.length < 13) state.pairs.push([pendingSocket, letter]);
    pendingSocket = null;
  } else {
    pendingSocket = letter;
  }
  renderPlugboard();
  update();
}

function renderPlugboard() {
  for (const letter of ALPHABET) {
    const socket = $(`socket-${letter}`);
    const pairIndex = state.pairs.findIndex((p) => p.includes(letter));
    const holes = socket.querySelector(".holes");

    socket.classList.toggle("paired", pairIndex !== -1);
    socket.classList.toggle("selected", pendingSocket === letter);

    if (pairIndex !== -1) {
      const partner = state.pairs[pairIndex].find((l) => l !== letter);
      holes.textContent = partner;
      socket.setAttribute("aria-label", `${letter}, connected to ${partner}. Click to unplug.`);
    } else {
      holes.textContent = "";
      socket.setAttribute(
        "aria-label",
        pendingSocket === letter ? `${letter}, selected. Click another letter to connect.` : `${letter}, not connected`,
      );
    }
  }

  const count = state.pairs.length;
  $("pair-count").textContent = `${count} cable${count === 1 ? "" : "s"}`;
  const input = $("plug-input");
  if (document.activeElement !== input) input.value = state.pairs.map((p) => p.join("")).join(" ");
  input.removeAttribute("aria-invalid");
}

$("plug-input").addEventListener("input", (e) => {
  const letters = e.target.value.toUpperCase().replace(/[^A-Z]/g, "");
  const unique = new Set(letters).size === letters.length;
  if (letters.length % 2 !== 0 || !unique || letters.length > 26) {
    e.target.setAttribute("aria-invalid", "true");
    return;
  }
  state.pairs = letters.match(/../g)?.map((p) => [p[0], p[1]]) ?? [];
  pendingSocket = null;
  renderPlugboard();
  update();
});

$("plug-input").addEventListener("blur", () => renderPlugboard());

$("plug-clear").addEventListener("click", () => {
  state.pairs = [];
  pendingSocket = null;
  renderPlugboard();
  update();
});

// ---------------------------------------------------------------------------------------------------------------------
// Signal path diagram
// ---------------------------------------------------------------------------------------------------------------------

const SVG_NS = "http://www.w3.org/2000/svg";

function svg(tag, attrs = {}, text) {
  const node = document.createElementNS(SVG_NS, tag);
  for (const [key, value] of Object.entries(attrs)) {
    // Presentation attributes can't use CSS variables, but the style property can.
    if (typeof value === "string" && value.startsWith("var(")) node.style.setProperty(key, value);
    else node.setAttribute(key, value);
  }
  if (text !== undefined) node.textContent = text;
  return node;
}

function renderPath(result) {
  const model = MODELS[state.model];
  const ROW = 15;
  const TOP = 50;
  const BLOCK = 54;
  const GAP = 42;
  const LEFT = 14;

  // Blocks are drawn in the order they sit in the machine, from the reflector on the left to the plugboard on the right.
  const blocks = [{ id: "reflector", name: REFLECTOR_LABELS[state.reflector].replace("UKW-", "UKW "), sub: "Reflector" }];
  if (model.fourth) {
    blocks.push({ id: "fourth", name: FOURTH_ROTORS[state.fourthRotor], sub: "Fourth", window: result?.fourthPosition });
  }
  SLOTS.forEach((slot, i) =>
    blocks.push({ id: slot.toLowerCase(), name: state.rotors[i], sub: slot, window: result?.positions?.[i] }),
  );
  blocks.push({ id: "plugboard", name: "Plugs", sub: "Plugboard" });
  blocks.forEach((b, i) => {
    b.left = LEFT + i * (BLOCK + GAP);
    b.right = b.left + BLOCK;
  });

  const keysX = blocks.at(-1).right + GAP;
  const width = keysX + 30;
  const height = TOP + 26 * ROW + 22;
  const y = (letter) => TOP + ALPHABET.indexOf(letter) * ROW + ROW / 2;
  const byId = Object.fromEntries(blocks.map((b) => [b.id, b]));

  const root = $("path-svg");
  root.setAttribute("viewBox", `0 0 ${width} ${height}`);
  const nodes = [];
  const muted = "var(--muted)";

  for (const b of blocks) {
    nodes.push(svg("rect", { x: b.left, y: TOP - 6, width: BLOCK, height: 26 * ROW + 12, rx: 7, fill: "none", stroke: "var(--line)" }));
    nodes.push(svg("text", { x: b.left + BLOCK / 2, y: 18, "text-anchor": "middle", "font-family": "var(--mono)", "font-weight": 700, "font-size": 13, fill: "var(--text)" }, b.name));
    nodes.push(svg("text", { x: b.left + BLOCK / 2, y: 33, "text-anchor": "middle", "font-size": 10, fill: muted }, b.sub));
    if (b.window) {
      nodes.push(svg("text", { x: b.left + BLOCK / 2, y: height - 2, "text-anchor": "middle", "font-family": "var(--mono)", "font-size": 10, fill: muted }, `at ${b.window}`));
    }
    for (let i = 0; i < 26; i++) {
      const cy = TOP + i * ROW + ROW / 2;
      if (b.id !== "reflector") nodes.push(svg("circle", { cx: b.left, cy, r: 1.6, fill: "var(--line)" }));
      nodes.push(svg("circle", { cx: b.right, cy, r: 1.6, fill: "var(--line)" }));
    }
  }

  const trace = result?.trace ?? [];
  const key = trace[0]?.in;
  const lamp = trace.at(-1)?.out;

  nodes.push(svg("text", { x: keysX, y: 18, "text-anchor": "middle", "font-weight": 700, "font-size": 11, fill: "var(--text)" }, "Keys"));
  nodes.push(svg("text", { x: keysX, y: 33, "text-anchor": "middle", "font-size": 10, fill: muted }, "& lamps"));
  for (const letter of ALPHABET) {
    const isKey = letter === key;
    const isLamp = letter === lamp;
    if (isKey || isLamp) {
      nodes.push(svg("circle", { cx: keysX, cy: y(letter), r: 7, fill: isKey ? "var(--text)" : "var(--accent)" }));
    }
    nodes.push(
      svg("text", {
        x: keysX,
        y: y(letter) + 3.5,
        "text-anchor": "middle",
        "font-family": "var(--mono)",
        "font-size": 10,
        "font-weight": isKey || isLamp ? 700 : 400,
        fill: isKey || isLamp ? "var(--bg)" : muted,
      }, letter),
    );
  }

  if (trace.length) {
    const reflectorIndex = trace.findIndex((s) => s.component === "reflector");
    const reflector = byId.reflector;
    const turnX = reflector.left + BLOCK * 0.4;
    const forward = [[keysX - 9, y(key)]];
    const back = [];
    const labels = [];

    trace.forEach((s, i) => {
      if (i === reflectorIndex) {
        forward.push([reflector.right, y(s.in)], [turnX, y(s.in)]);
        back.push([turnX, y(s.in)], [turnX, y(s.out)], [reflector.right, y(s.out)]);
        labels.push({ x: reflector.right + GAP / 2, y: y(s.out) + 11, text: s.out, color: "var(--accent)" });
        return;
      }
      const b = byId[s.component];
      if (i < reflectorIndex) {
        forward.push([b.right, y(s.in)], [b.left, y(s.out)]);
        labels.push({ x: b.left - GAP / 2, y: y(s.out) - 4, text: s.out, color: "var(--text)" });
      } else {
        back.push([b.left, y(s.in)], [b.right, y(s.out)]);
        labels.push({ x: b.right + GAP / 2, y: y(s.out) + 11, text: s.out, color: "var(--accent)" });
      }
    });
    back.push([keysX - 9, y(lamp)]);

    const line = (points, color) =>
      svg("polyline", {
        points: points.map((p) => p.join(",")).join(" "),
        fill: "none",
        stroke: color,
        "stroke-width": 1.75,
        "stroke-linejoin": "round",
        "stroke-linecap": "round",
      });
    nodes.push(line(back, "var(--accent)"), line(forward, "var(--text)"));

    for (const l of labels) {
      nodes.push(
        svg("text", {
          x: l.x,
          y: l.y,
          "text-anchor": "middle",
          "font-family": "var(--mono)",
          "font-size": 10,
          "font-weight": 700,
          fill: l.color,
          stroke: "var(--bg)",
          "stroke-width": 3,
          "paint-order": "stroke",
        }, l.text),
      );
    }
  }

  root.replaceChildren(...nodes);
  renderPathSummary(trace, byId);
}

function renderPathSummary(trace, byId) {
  const summary = $("path-summary");
  if (!trace.length) {
    summary.textContent = "Press a key to see the path it takes through the machine.";
    return;
  }

  const name = (component) => {
    if (component === "plugboard") return "plugboard";
    if (component === "reflector") return "reflector";
    return `rotor ${byId[component].name}`;
  };

  const parts = [el("b", { text: `Key ${trace[0].in}` })];
  for (const s of trace) parts.push(` → ${name(s.component)} → ${s.out}`);
  parts.push(" → ", el("b", { text: `lamp ${trace.at(-1).out}` }));
  summary.replaceChildren(...parts);
}

// ---------------------------------------------------------------------------------------------------------------------
// Running the engine
// ---------------------------------------------------------------------------------------------------------------------

function setWindow(id, letter) {
  const node = $(id);
  if (!node || node.textContent === letter) return;
  node.textContent = letter;
}

function update() {
  saveState();
  if (!engineReady) return;

  const result = JSON.parse(window.ultraRun(JSON.stringify(engineConfig()), state.text));
  const error = $("error");

  if (result.error) {
    error.textContent = result.error;
    error.hidden = false;
    return;
  }
  error.hidden = true;
  lastResult = result;

  const count = result.input.length;
  $("letter-count").textContent = `${count} letter${count === 1 ? "" : "s"}`;
  $("output-text").textContent = groupsOfFive(result.output);

  [...result.positions].forEach((letter, i) => setWindow(`window-${i}`, letter));
  if (result.fourthPosition) setWindow("window-fourth", result.fourthPosition);

  const lit = result.output.at(-1);
  for (const letter of ALPHABET) $(`lamp-${letter}`).classList.toggle("lit", letter === lit);

  renderPath(result);
}

async function loadEngine() {
  const status = $("engine-status");
  try {
    const go = new Go();
    const ready = new Promise((resolve) => (window.onUltraReady = resolve));

    let instance;
    try {
      ({ instance } = await WebAssembly.instantiateStreaming(fetch("ultra.wasm"), go.importObject));
    } catch {
      // Some servers don't send the application/wasm content type that streaming needs.
      const response = await fetch("ultra.wasm");
      if (!response.ok) throw new Error(`ultra.wasm: HTTP ${response.status}`);
      ({ instance } = await WebAssembly.instantiate(await response.arrayBuffer(), go.importObject));
    }

    go.run(instance);
    await ready;

    engineReady = true;
    status.dataset.state = "ready";
    status.textContent = "";
    update();
  } catch (err) {
    console.error(err);
    status.dataset.state = "error";
    status.textContent = "Engine failed to load";
    const error = $("error");
    error.textContent = `The simulator couldn't start: ${err.message}`;
    error.hidden = false;
  }
}

// ---------------------------------------------------------------------------------------------------------------------
// Start up
// ---------------------------------------------------------------------------------------------------------------------

buildBoards();
buildPlugboard();
renderSettings();
renderPlugboard();
renderPath(null);
$("input-text").value = state.text;
loadEngine();
