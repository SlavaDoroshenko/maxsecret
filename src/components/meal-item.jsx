import {
  Panel,
  Grid,
  Container,
  Flex,
  Typography,
  Button,
  CellAction,
  CellList,
  CellSimple,
} from "@maxhub/max-ui";

import Plus from "../assets/plus.svg?react";
import BreakFastIcon from "../assets/breakfast.svg?react";

const MealItem = ({ id, label, current, max }) => (
  <CellList
    mode="full-width"
    filled={true}
    className="shadow-xl rounded-xl items-center overflow-hidden h-20 mb-3"
  >
    <CellAction
      before={<BreakFastIcon className="w-6 h-6" />}
      height="compact"
      mode="primary"
      onClick={() => {}}
      showChevron={false}
      className="h-20"
    >
      <Typography.Body className="flex justify-between items-center text-lg">
        <Typography.Body className="p-0 font-bold text-2xl">
          {" "}
          {label}
        </Typography.Body>
        <div className="flex flex-col items-center">
          <Typography.Body className="p-0 font-bold">
            {current} ккал
          </Typography.Body>
          <Typography.Body className="p-0">из {max} ккал</Typography.Body>
        </div>
        <Plus className="w-4 h-4" />
      </Typography.Body>
    </CellAction>
  </CellList>
);

export default MealItem;
