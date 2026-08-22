# langmuir-iso

A self-contained toolkit for the **Langmuir monolayer adsorption isotherm**. Given
an adsorbate, a temperature, and a Langmuir constant `K` (or an isotherm
dataset), it computes surface coverage `θ`, adsorbed amount `q`, the
`isotherm` curve `θ(p)`, the Henry-law limit, and the inverse problem of
recovering pressure from a target coverage — for both single-component and
multi-component competitive adsorption.

This is an adsorption *kernel*: it models equilibrium gas–solid uptake at a
single interface, not a full reactor or process simulator.

## What it computes

The Langmuir isotherm assumes a single homogeneous surface with one site per
adsorbate:

```
theta(p)    = K·p / (1 + K·p)            # fractional coverage
q(p)        = q_max · theta(p)           # adsorbed amount
theta ~ K·p                               # Henry-law limit as p -> 0
p(theta)    = theta / (K·(1 - theta))    # inverse: pressure from coverage
```

For binary competitive adsorption (`mode = "competitive"`), each component `i`
follows the multicomponent Langmuir form:

```
theta_i = K_i·p_i / (1 + Σ_j K_j·p_j)
```

## Inputs, outputs, and boundaries

- `K` must be non-negative; `q_max` must be positive.
- Pressure `p` must be ≥ 0.
- At zero pressure coverage is exactly zero; as `p -> ∞` coverage tends to 1.
- For competitive mode at least two components must be supplied.
- Invalid input (negative pressure, missing `q_max`, etc.) returns a clear error
  instead of a wrong number.

## Usage

Start the web server and JSON API:

```bash
go run . -http :8080
```

Then open <http://localhost:8080/>. The page can load the bundled example
(`example/n2-77k.json`, nitrogen at 77 K) and draw the `q–p` isotherm fetched
from the backend.

Run the offline CLI over the bundled example:

```bash
go run . compute
# N2 at 77 K, driver=single, unit=mmol/g
# K=..., q_max=...
# pressure       theta         q
# ...
```

And a couple of API calls:

```bash
curl -s -X POST localhost:8080/api/iso   -d '{"k":0.5,"pmax":1.0,"qmax":2.0,"p":0.3}'
curl -s -X POST localhost:8080/api/curve -d '{"k":0.5,"pmax":1.0,"qmax":2.0,"steps":20}'
curl -s -X POST localhost:8080/api/inverse -d '{"k":0.5,"theta":0.5}'
```

## API endpoints

| Method | Path           | Body                          | Returns                              |
|--------|----------------|-------------------------------|--------------------------------------|
| POST   | `/api/iso`     | `{k,qmax,p,pmax}`             | `{theta,q}`                           |
| POST   | `/api/curve`   | `{k,qmax,pmax,steps}`         | `{points:[{p,theta,q}...]}`          |
| POST   | `/api/inverse` | `{k,theta}`                   | `{p}`                                |
| POST   | `/api/comp`    | `{components:[{k,p}], qmax}`   | `{components:[{theta,q}...]}`         |
| GET    | `/api/example` | —                             | the bundled `n2-77k.json`            |

## Build and test

```bash
go build ./...
go test ./...
```
