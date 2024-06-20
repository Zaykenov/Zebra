import React from "react";

export interface CheckboxNode {
  label: string;
  value: string;
  isSelected: boolean;
  isDisabled: boolean;
  isCollapsed: boolean;
  children: CheckboxNode[];
}

export const traverseAndToggleNode = (
  data: CheckboxNode[],
  value: string,
  key: "isSelected" | "isDisabled" | "isCollapsed" = "isSelected"
) => {
  let toggled = false;

  if (data) {
    data.forEach((node) => {
      if (node.value === value) {
        node[key] = !node[key];
        toggled = true;
      }
    });

    if (!toggled) {
      for (let i = 0; i < data.length; i++) {
        const node = data[i];
        const { data: newChildren, toggled } = traverseAndToggleNode(
          node.children,
          value
        );
        if (toggled) {
          node.children = newChildren;
          break;
        }
      }
    }
  }
  return { data, toggled };
};

export const toggleNodeSelected = (value: string, setData: any): void => {
  setData((data: CheckboxNode) => {
    const { data: newChildren } = traverseAndToggleNode(data.children, value);
    return { ...data, children: newChildren };
  });
};

export const setAllChildren = (node: CheckboxNode, value: boolean) => {
  console.log(value);
  let flat: CheckboxNode[] = getFlattedChildren(node);

  flat.forEach((child: CheckboxNode) => {
    child.isSelected = value;
  });
};

export const traverseAndToggleNodeChildren = (
  data: CheckboxNode[],
  value: string
) => {
  let toggled = false;
  data.forEach((node) => {
    if (node.value === value) {
      toggled = true;
      let newValue = true;
      // If all selected, change them to false
      if (isAllSelected(node)) {
        newValue = false;
      }

      node.isSelected = newValue;
      setAllChildren(node, newValue);
    }
  });

  if (!toggled) {
    for (let i = 0; i < data.length; i++) {
      let node = data[i];
      const { data: newChildren, toggled } = traverseAndToggleNodeChildren(
        node.children,
        value
      );
      node.children = newChildren;
      if (toggled) {
        break;
      }
    }
  }

  return { data, toggled };
};

export const toggleAllChildren = (value: string, setData: any) => {
  console.log("toggleAllChildren");
  setData((data: CheckboxNode) => {
    const { data: newChildren } = traverseAndToggleNodeChildren(
      data.children,
      value
    );
    return { ...data, children: newChildren };
  });
};

export const CheckItem = ({
  node,
  setData,
}: {
  node: CheckboxNode;
  setData: (node: CheckboxNode) => void;
}) => {
  return (
    <div style={{ width: "100%" }}>
      <div
        style={{
          display: "flex",
          justifyContent: "space-between",
          width: "100%",
          // border: "1px solid red"
        }}
      >
        <label style={{ color: node.isDisabled ? "gray" : "black" }}>
          <input
            type="checkbox"
            checked={node.isSelected}
            disabled={node.isDisabled}
            onChange={() => {
              toggleNodeSelected(node.value, setData);
            }}
          />
          <span>{node.label}</span>
        </label>
        <div>
          <input
            type="checkbox"
            disabled={node.isDisabled}
            onChange={() => {
              toggleAllChildren(node.value, setData);
            }}
            checked={isAllSelected(node)}
          />
        </div>
      </div>
      <div style={{ paddingLeft: "16px", paddingTop: "8px" }}>
        <NestedChecklist data={node.children} setData={setData} />
      </div>
    </div>
  );
};

export const NestedChecklist = ({
  data,
  setData,
}: {
  data: CheckboxNode[];
  setData: (node: CheckboxNode) => void;
}) => {
  return (
    <>
      {data.map((node) => (
        <CheckItem node={node} setData={setData} />
      ))}
    </>
  );
};

export const isAllSelected = (node: CheckboxNode) => {
  let allSelected = true;

  let flat: CheckboxNode[] = getFlattedChildren(node);

  flat.forEach((child: CheckboxNode) => {
    if (!child.isSelected) {
      allSelected = false;
    }
  });

  return allSelected;
};

export const getFlattedChildren = (node: CheckboxNode): CheckboxNode[] => {
  let flat = [node];
  node.children.forEach((child) => {
    flat = [...flat, ...getFlattedChildren(child)];
  });
  return flat;
};
