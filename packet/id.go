package packet

import "fmt"

type PacketID byte

// Client -> Server
const (
	IDClientKeepAlive PacketID = iota
	IDClientChatMessage
	IDClientUseEntity
	IDClientPlayer
	IDClientPlayerPosition
	IDClientPlayerLook
	IDClientPlayerPositionAndLook
	IDClientPlayerDigging
	IDClientPlayerBlockPlacement
	IDClientHeldItemChange
	IDClientAnimation
	IDClientEntityAction
	IDClientSteerVehicle
	IDClientCloseWindow
	IDClientClickWindow
	IDClientConfirmTransaction
	IDClientCreativeInventoryAction
	IDClientEnchantItem
	IDClientUpdateSign
	IDClientPlayerAbilities
	IDClientTabComplete
	IDClientClientSettings
	IDClientClientStatus
	IDClientPluginMessage
)

// Server -> Client
const (
	IDServerKeepAlive PacketID = iota
	IDServerJoinGame
	IDServerChatMessage
	IDServerTimeUpdate
	IDServerEntityEquipment
	IDServerSpawnPosition
	IDServerUpdateHealth
	IDServerRespawn
	IDServerPlayPositionAndLook
	IDServerHeldItemChange
	IDServerUseBed
	IDServerAnimation
	IDServerSpawnPlayer
	IDServerCollectItem
	IDServerSpawnObject
	IDServerSpawnMob
	IDServerSpawnPainting
	IDServerExperienceOrb
	IDServerEntityVelocity
	IDServerDestroyEntities
	IDServerEntity
	IDServerEntityRelativeMove
	IDServerEntityLook
	IDServerEntityLookAndRelativeMove
	IDServerEntityTeleport
	IDServerEntityHeadLook
	IDServerEntityStatus
	IDServerAttachEntity
	IDServerEntityMetadata
	IDServerEntityEffect
	IDServerRemoveEntityEffect
	IDServerSetExperience
	IDServerEntityProperties
	IDServerChunkData
	IDServerMultiBlockChange
	IDServerBlockChange
	IDServerBlockAction
	IDServerBlockBreakAnimation
	IDServerMapChunkBulk
	IDServerExplosion
	IDServerEffect
	IDServerSoundEffect
	IDServerParticle
	IDServerChangeGameState
	IDServerSpawnGlobalEntity
	IDServerOpenWindow
	IDServerCloseWindow
	IDServerSetSlot
	IDServerWindowItems
	IDServerWindowProperty
	IDServerConfirmTransaction
	IDServerUpdateSign
	IDServerMaps
	IDServerUpdateBlockEntity
	IDServerSignEditorOpen
	IDServerStatistics
	IDServerPlayListItems
	IDServerPlayerAbilities
	IDServerTabComplete
	IDServerScoreboardObjective
	IDServerUpdateScore
	IDServerDisplayScoreboard
	IDServerTeams
	IDServerPluginMessage
	IDServerDisconnect
)

func (pID PacketID) String() string {
	return fmt.Sprintf("0x%x", byte(pID))
}
