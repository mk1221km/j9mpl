(* Data-oriented ring buffer in OCaml.
   No structs, no objects, no type system nonsense.
   Flat array, scalar head/count, freestanding functions. *)

let capacity = 1024

(* ringPush: store value at head, return new head and count *)
let ring_push buf head count value =
  buf.(head) <- value;
  let head = (head + 1) mod capacity in
  let count = if count < capacity then count + 1 else count in
  (head, count)

(* ringAvg: arithmetic mean of stored elements. 0 if empty. *)
let ring_avg buf count =
  if count = 0 then 0.0
  else
    let rec sum i acc =
      if i >= count then acc
      else sum (i + 1) (acc +. buf.(i))
    in
    sum 0 0.0 /. (float_of_int count)

(* ringReadRange: copy up to n most recent elements into a new array.
   Index math: (head - n + i + capacity) % capacity (0-indexed, same as Go). *)
let ring_read_range buf head count n =
  if n <= 0 || count = 0 then [||]
  else
    let n = if n > count then count else n in
    Array.init n (fun i ->
      let idx = (head - n + i + capacity) mod capacity in
      buf.(idx))

(* ringUtilization: ratio of stored count to capacity. *)
let ring_utilization count =
  (float_of_int count) /. (float_of_int capacity)

(* Filter empty strings from a list (handles multiple spaces). *)
let rec filter_empty = function
  | [] -> []
| "" :: rest -> filter_empty rest
| s :: rest -> s :: filter_empty rest

(* Split a trimmed line into command words. *)
let split_words line =
  filter_empty (String.split_on_char ' ' line)

(* Main: read stdin line by line, dispatch commands. *)
let () =
  let storage = Array.make capacity 0.0 in
  let rec loop head count =
    match try Some (input_line stdin) with End_of_file -> None with
    | None -> ()
    | Some line ->
      let trimmed = String.trim line in
      if trimmed = "" then loop head count
      else
        let parts = split_words trimmed in
        (match parts with
        | "push" :: value_str :: _ ->
          let value = float_of_string value_str in
          let (head', count') = ring_push storage head count value in
          loop head' count'
        | "avg" :: _ ->
          Printf.printf "%.1f\n" (ring_avg storage count);
          loop head count
        | "readRange" :: rest ->
          let n = match rest with
            | n_str :: _ -> (try int_of_string n_str with Failure _ -> 1)
            | [] -> 1
          in
          let dest = ring_read_range storage head count n in
          Array.iter (fun v -> Printf.printf "%.1f\n" v) dest;
          loop head count
        | "utilization" :: _ ->
          Printf.printf "%.4f\n" (ring_utilization count);
          loop head count
        | _ -> loop head count)
  in
  loop 0 0;
  flush stdout
