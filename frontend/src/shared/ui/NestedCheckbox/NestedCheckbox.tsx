import React, {
  Dispatch,
  FC,
  SetStateAction,
  useEffect,
  useState,
} from "react";
import CheckboxTree, { Node, OnCheckNode } from "react-checkbox-tree";
import { ChevronDownIcon, ChevronRightIcon } from "@heroicons/react/24/outline";
import { getAllValuesFromNodes } from "@modules/InventoryForm/InventoryForm";

export interface NestedCheckboxProps {
  nodes: Node[];
  searchedNodes: Node[];
  checked: string[];
  expanded: string[];
  setChecked: Dispatch<SetStateAction<string[]>>;
}

const NestedCheckbox: FC<NestedCheckboxProps> = ({
  nodes,
  searchedNodes,
  checked,
  expanded: defaultExpanded,
  setChecked,
}) => {
  const [expanded, setExpanded] = useState<string[]>([]);

  useEffect(() => {
    setExpanded(defaultExpanded);
  }, [defaultExpanded]);

  // const addNewCheck = (newChecked: string[]) => {
  //   const searchedNodesValues = getAllValuesFromNodes(searchedNodes, true);
  //   const allNodesValues = getAllValuesFromNodes(nodes, true);
  //   searchedNodesValues.length !== allNodesValues.length
  //     ? setChecked((prevState) => {
  //         let arr: string[] = Array.from(new Set(prevState.concat(newChecked)));
  //         const isRemove = newChecked.every((elem) => prevState.includes(elem));
  //         arr = isRemove
  //           ? arr.filter(
  //               (elem) =>
  //                 !searchedNodesValues.includes(elem) ||
  //                 newChecked.includes(elem)
  //             )
  //           : arr;
  //         return arr;
  //       })
  //     : setChecked(newChecked);
  // };

  const addNewCheck = (newChecked: string[]) => {
    setChecked((prevState) => {
      let arr: string[] = Array.from(new Set(prevState.concat(newChecked)));
      const isRemove = newChecked.every((elem) => prevState.includes(elem));
      arr = isRemove ? prevState.filter((elem) => newChecked.includes(elem)) : arr;
      return arr;
    });
  };  
  

  return (
    <CheckboxTree
      nodes={searchedNodes}
      checked={checked}
      expanded={expanded}
      onCheck={(newChecked) => {
        addNewCheck(newChecked);
      }}
      onExpand={(expanded) => setExpanded(expanded)}
      icons={{
        expandClose: <ChevronRightIcon className="w-4 h-4 inline" />,
        expandOpen: <ChevronDownIcon className="w-4 h-4 inline" />,
      }}
      expandOnClick
      onClick={(node: OnCheckNode) => {
        if (!node.children) {
          checked.includes(node.value)
            ? setChecked([
                ...checked.slice(0, checked.indexOf(node.value)),
                ...checked.slice(checked.indexOf(node.value) + 1),
              ])
            : addNewCheck([...checked, node.value]);
        }
      }}
    />
  );
};

export default NestedCheckbox;
